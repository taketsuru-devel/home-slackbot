package awssqsutil

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
)

type SQSReceiveCallback func(*types.Message) error

type SQS struct {
	queueUrl            *string
	sqsStr              *sqs.Client
	receiveMessageInput *sqs.ReceiveMessageInput
	callback            SQSReceiveCallback
	callbackTimeout     time.Duration
	wg                  *sync.WaitGroup
	ReceiveWaitSecond   int32
}

//  receive -> callback -> executeが成功したらdelete
//  まずは再送とかpush側とか考えない

const SQS_BATCH_MAX = 10

func InitSqs(sqsName string, regionName string, profileName string) *SQS {
	config := aws.Config{}

	if regionName != "" {
		config.Region = regionName
	}
	ret := &SQS{
		queueUrl: aws.String(sqsName),
		sqsStr:   sqs.NewFromConfig(config),
		wg:       &sync.WaitGroup{},
	}
	ret.receiveMessageInput = &sqs.ReceiveMessageInput{
		QueueUrl:            ret.queueUrl,
		MaxNumberOfMessages: SQS_BATCH_MAX,
		WaitTimeSeconds:     10, //暫定
	}
	return ret
}

//sdkのエラーハンドリングは後回し
func (s *SQS) Send(ctx context.Context, msgs []*string) {
	sendMsg := sqs.SendMessageBatchInput{
		Entries:  make([]types.SendMessageBatchRequestEntry, len(msgs)),
		QueueUrl: s.queueUrl,
	}
	for i, msg := range msgs {
		sendMsg.Entries[i] = types.SendMessageBatchRequestEntry{
			Id:          uuid.New().String(),
			MessageBody: msg,
		}
	}
	s.SendDetail(ctx, &sendMsg)
}

func (s *SQS) SendDetail(ctx context.Context, msgs *sqs.SendMessageBatchInput) {
	s.sqsStr.SendMessageBatch(ctx, msgs)
}

func (s *SQS) StartReceive(ctx context.Context, cb SQSReceiveCallback, cbTimeout time.Duration) {
	s.wg.Add(1)
	defer s.wg.Done()
	for {
		if received, err := s.sqsStr.ReceiveMessage(ctx, s.receiveMessageInput); err != nil {
			if ctx.Err() != nil {
				return
			}
		} else {
			deleteMsgInput := sqs.DeleteMessageBatchInput{
				Entries:  make([]types.DeleteMessageBatchRequestEntry, 0, len(received.Messages)),
				QueueUrl: s.queueUrl,
			}
			for _, msg := range received.Messages {
				if execErr := s.execCallback(ctx, cb, &msg, cbTimeout); execErr != nil {
					//ログ表示
				} else {
					deleteMsgInput.Entries = append(deleteMsgInput.Entries, types.DeleteMessageBatchRequestEntry{
						Id:            msg.MessageId,
						ReceiptHandle: msg.ReceiptHandle,
					})
				}
			}
			if len(deleteMsgInput.Entries) > 0 {
				go s.sqsStr.DeleteMessageBatch(ctx, &deleteMsgInput)
			}
		}
	}
}

func (s *SQS) Close() {
	s.wg.Wait()
}

func (s *SQS) execCallback(ctx context.Context, cb SQSReceiveCallback, msg *types.Message, cbTimeout time.Duration) (retErr error) {
	s.wg.Add(1)
	defer s.wg.Done()
	thisCtx, cancel := context.WithTimeout(ctx, cbTimeout)
	defer cancel()
	cbRetCh := make(chan error, 1) //sync poolか？
	go func() {
		//ここにwg.Done()を置くべきか
		defer close(cbRetCh)
		cbRetCh <- cb(msg)
	}()
	select {
	case err := <-cbRetCh:
		retErr = err
	case <-thisCtx.Done():
		retErr = thisCtx.Err()
	}
	return
}
