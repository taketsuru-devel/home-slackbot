package awssqsutil

import (
	"context"
	//"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
)

type testReceiveStr struct {
	t        *testing.T
	expected []*string
	cnt      int
}

func (ts *testReceiveStr) receiveCbImpl(msg *types.Message) error {
	ts.t.Logf("received: " + *msg.Body)
	assert.Equal(ts.t, *ts.expected[ts.cnt], *msg.Body)
	ts.cnt++
	return nil
}

//  send -> receive & deleteの確認
func TestSQS(t *testing.T) {
	t1 := "test1"
	t2 := "test2"
	testMsg := []*string{&t1, &t2}

	ts := testReceiveStr{t: t, expected: testMsg}
	input := SQSClientInput{
		QueueUrl:               os.Getenv("TEST_FIFO_SQS"),
		ReceiveCallback:        ts.receiveCbImpl,
		ReceiveCallbackTimeout: time.Duration(100 * time.Millisecond),
	}

	sqsStr := GetSqsClient(&input)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go sqsStr.StartReceive(ctx)
	t.Log("receive start")

	err := sqsStr.Send(context.Background(), testMsg)
	if err != nil {
		t.Error(err)
	}
	t.Log("send done")

	time.Sleep(1 * time.Second)
	t.Log("timeout")

	cancel()
	sqsStr.Close()
	assert.Equal(t, len(testMsg), ts.cnt)
}
