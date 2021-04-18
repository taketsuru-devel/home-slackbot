package util

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func IotInvoke(target string, command string) error {
	sess, _ := session.NewSessionWithOptions(session.Options{
		//Profile; "default",
		Config: aws.Config{
			Region:                        aws.String("us-west-2"),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	})
	svc := lambda.New(sess)
	input := &lambda.InvokeInput{
		FunctionName: aws.String("home-iot-invoker"),
		Payload:      []byte(fmt.Sprintf("{\"target\":\"%s\", \"command\":\"%s\"}", target, command)),
		//Qualifier:    aws.String("1"),
	}

	if resp, err := svc.Invoke(input); err != nil {
		return err
	} else if *resp.StatusCode != int64(200) {
		return fmt.Errorf("%v", resp)
	}
	return nil
}
