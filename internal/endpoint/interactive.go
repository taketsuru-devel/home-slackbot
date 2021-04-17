package endpoint

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/slack-go/slack"
	"net/http"
	"os"
	"strings"
)

type InteractiveEndpoint struct{}

func (i *InteractiveEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	// slack.New("YOUR_TOKEN_HERE", slack.OptionDebug(true))
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"), slack.OptionDebug(true))
	fmt.Println("interactive")
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v", err)
		//ここでエラーだとChannelの取得もできない
		//api.PostMessage(payload.Channel, slack.MsgOptionText(fmt.Printf("指令の解析に失敗しました: %v", err), false))
		return
	}
	channelId := payload.Channel.GroupConversation.Conversation.ID
	fmt.Printf("%#v\n", payload.ActionCallback.BlockActions[0].Value)
	api.PostMessage(channelId, slack.MsgOptionText(fmt.Sprintf("%sを受け付けました", payload.ActionCallback.BlockActions[0].Text.Text), false))
	sess, _ := session.NewSessionWithOptions(session.Options{
		//Profile; "default",
		Config: aws.Config{
			Region:                        aws.String("us-west-2"),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	})
	svc := lambda.New(sess)
	command := payload.ActionCallback.BlockActions[0].Value
	commands := strings.Split(command, ":")
	input := &lambda.InvokeInput{
		FunctionName: aws.String("home-iot-invoker"),
		Payload:      []byte(fmt.Sprintf("{\"target\":\"%s\", \"command\":\"%s\"}", commands[0], commands[1])),
		//Qualifier:    aws.String("1"),
	}
	if resp, err := svc.Invoke(input); err != nil {
		api.PostMessage(channelId, slack.MsgOptionText(fmt.Sprintf("処置に失敗しました: %v", err), false))
	} else if *resp.StatusCode != int64(200) {
		api.PostMessage(channelId, slack.MsgOptionText(fmt.Sprintf("処置に失敗しました: %v", resp), false))
	} else {
		api.PostMessage(channelId, slack.MsgOptionText("処置しました", false))
	}
}
