package awssqsutil

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
)

type SQSReceiveCallback func(*types.Message) error

type SQS struct {
	queueUrl                    *string
	sqsStr                      *sqs.Client
	receiveMessageInput         *sqs.ReceiveMessageInput
	callback                    SQSReceiveCallback
	callbackTimeout             time.Duration
	wg                          *sync.WaitGroup
	isFifo                      bool
	isContentBasedDeduplication bool
}

type SQSClientInput struct {
	QueueUrl                    string
	Region                      string
	Profile                     string
	ReceiveCallback             SQSReceiveCallback
	ReceiveCallbackTimeout      time.Duration
	ReceiveMaxNumberOfMessages  int32
	ReceiveWaitSecond           int32
	IsContentBasedDeduplication bool
}

//  receive -> callback -> executeが成功したらdelete
//  まずは再送とかpush側とか考えない

const RECEIVE_MAX_NUMBER_OF_MESSAGES_DEFAULT = 10
const RECEIVE_WAIT_TIME_SECONDS_DEFAULT = 15

func GetSqsClient(input *SQSClientInput) *SQS {
	cfgOpts := make([]func(*config.LoadOptions) error, 0, 5)
	if input.Region != "" {
		cfgOpts = append(cfgOpts, config.WithRegion(input.Region))
	}
	if input.Profile != "" {
		cfgOpts = append(cfgOpts, config.WithSharedConfigProfile(input.Profile))
	}
	cfg, _ := config.LoadDefaultConfig(context.TODO(), cfgOpts...)

	isFifo, _ := regexp.MatchString("\\.fifo$", input.QueueUrl)
	ret := &SQS{
		queueUrl:                    aws.String(input.QueueUrl),
		sqsStr:                      sqs.NewFromConfig(cfg),
		wg:                          &sync.WaitGroup{},
		callback:                    input.ReceiveCallback,
		callbackTimeout:             input.ReceiveCallbackTimeout,
		isFifo:                      isFifo,
		isContentBasedDeduplication: input.IsContentBasedDeduplication,
	}
	ret.receiveMessageInput = &sqs.ReceiveMessageInput{
		QueueUrl:            ret.queueUrl,
		MaxNumberOfMessages: input.ReceiveMaxNumberOfMessages,
		WaitTimeSeconds:     input.ReceiveWaitSecond,
	}
	//未定義 = デフォ値の処理であり上限下限判定はsdkコール時任せ
	if ret.receiveMessageInput.MaxNumberOfMessages == 0 {
		ret.receiveMessageInput.MaxNumberOfMessages = RECEIVE_MAX_NUMBER_OF_MESSAGES_DEFAULT
	}
	if ret.receiveMessageInput.WaitTimeSeconds == 0 {
		ret.receiveMessageInput.WaitTimeSeconds = RECEIVE_WAIT_TIME_SECONDS_DEFAULT
	}
	return ret
}

//sdkのエラーハンドリングは後回し
func (s *SQS) Send(ctx context.Context, msgs []*string) error {
	sendMsg := sqs.SendMessageBatchInput{
		Entries:  make([]types.SendMessageBatchRequestEntry, len(msgs)),
		QueueUrl: s.queueUrl,
	}
	var msgGroupId *string
	if s.isFifo {
		msgGroupId = aws.String(uuid.New().String())
	}
	for i, msg := range msgs {
		var deduplicationId *string
		if s.isFifo && !s.isContentBasedDeduplication {
			deduplicationId = aws.String(uuid.New().String())
		}
		sendMsg.Entries[i] = types.SendMessageBatchRequestEntry{
			Id:                     aws.String(uuid.New().String()),
			MessageGroupId:         msgGroupId,
			MessageDeduplicationId: deduplicationId,
			MessageBody:            msg,
		}
	}
	return s.SendDetail(ctx, &sendMsg)
}

//再送処理は後回し
func (s *SQS) SendDetail(ctx context.Context, msgs *sqs.SendMessageBatchInput) error {
	_, err := s.sqsStr.SendMessageBatch(ctx, msgs)
	//fmt.Printf("success: %#v\n", ret.Successful)
	//fmt.Printf("failed : %#v\n", ret.Failed)
	return err
}

func (s *SQS) StartReceive(ctx context.Context) {
	s.wg.Add(1)
	defer s.wg.Done()
	for {
		if received, err := s.sqsStr.ReceiveMessage(ctx, s.receiveMessageInput); err != nil {
			if ctx.Err() != nil {
				return
			}
		} else if len(received.Messages) > 0 {
			deleteMsgInput := sqs.DeleteMessageBatchInput{
				Entries:  make([]types.DeleteMessageBatchRequestEntry, 0, len(received.Messages)),
				QueueUrl: s.queueUrl,
			}
			for _, msg := range received.Messages {
				if execErr := s.execCallback(ctx, s.callback, &msg, s.callbackTimeout); execErr != nil {
					//ログ表示
					fmt.Println(execErr)
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
