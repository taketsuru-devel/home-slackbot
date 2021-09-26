package endpoint

import (
	"fmt"
	"github.com/followedwind/slackbot/internal/interactive"
	"github.com/followedwind/slackbot/internal/util"
	"github.com/slack-go/slack"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/taketsuru-devel/gorilla-microservice-skeleton/slackwrap"
)

func GetInteractiveHandler() *slackwrap.InteractiveEndpoint {
	return &slackwrap.InteractiveEndpoint{
		Handler: &interactiveHandler{},
	}
}

type interactiveHandler struct{}

func (ih *interactiveHandler) Handle() slackwrap.InteractiveHandlerFunc {
	return (func(w http.ResponseWriter, r *http.Request, ic *slack.InteractionCallback) {
		channelId := ic.Channel.GroupConversation.Conversation.ID
		command := ic.ActionCallback.BlockActions[0].Value
		commandDisp := ic.ActionCallback.BlockActions[0].Text.Text
		OriginTs := ic.MessageTs

		util.DebugLog(fmt.Sprintf("ic:%#v", ic), 0)
		util.DebugLog(fmt.Sprintf("command:%#v", ic.ActionCallback.BlockActions[0].Value), 0)
		//とりあえず返事
		api := util.GetSlackClient()

		if ic.BlockActionState == nil {
			util.DebugLog("未対応のパスです", 0)
		} else if len(ic.BlockActionState.Values) > 0 {
			util.DebugLog(fmt.Sprintf("block_actions:%#v", ic.BlockActionState), 0)
			//map[string]map[string]BlockAction
			//一次のキーが何を指してるか不明
			//二次のキーはevent_id
			for _, vmap := range ic.BlockActionState.Values {
				for eventId, state := range vmap {
					if eventId == "pwselect" {
						api.PostMessage(channelId, slack.MsgOptionText(fmt.Sprintf("%sを受け付けました", state.SelectedOption.Value), false))
						//Dynamo getAll and show list
						svc := dynamodb.New(session.New(), &aws.Config{Region: aws.String("ap-northeast-1")})
						input := &dynamodb.QueryInput{
							TableName:                 aws.String("PassDb"),
							ProjectionExpression:      aws.String("PassData"),
							KeyConditionExpression:    aws.String("Service = :v1"),
							ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":v1": &dynamodb.AttributeValue{S: aws.String(state.SelectedOption.Value)}},
						}
						if result, err := svc.Query(input); err != nil {
							fmt.Println(err)
						} else {
							//一覧表示してlistにしてイベントを作る
							for _, data := range result.Items {
								dataMap := data["PassData"].M
								strbuf := make([]string, 0, len(dataMap))
								for k, v := range dataMap {
									strbuf = append(strbuf, fmt.Sprintf("%s: %s", k, *v.S))
								}
								api := util.GetSlackClient()
								if _, ts, err := api.PostMessage(channelId, slack.MsgOptionText(strings.Join(strbuf, "\n"), false)); err != nil {
									util.ErrorLog(err.Error(), 0)
								} else {
									//連続だとエラーが出るので少し寝る
									time.Sleep(100 * time.Millisecond)
									if _, _, err := api.PostMessage(channelId, slack.MsgOptionText("60秒後に削除されます", false), slack.MsgOptionTS(ts)); err != nil {
										util.ErrorLog(err.Error(), 0)
									}
									go func() {
										time.Sleep(60 * time.Second)
										api.DeleteMessage(channelId, ts)
									}()
								}
							}
						}
						go func() {
							time.Sleep(1 * time.Second)
							util.InfoLog("delete test", 0)
							if _, _, err := api.DeleteMessage(channelId, OriginTs); err != nil {
								util.ErrorLog(err.Error(), 0)

							}
						}()
					}
				}
			}
		} else {

			//とりあえず返事
			api := util.GetSlackClient()
			api.PostMessage(channelId, slack.MsgOptionText(fmt.Sprintf("%sを受け付けました", commandDisp), false))

			//モノに指令
			commands := strings.Split(command, ":")
			responseText := "処置しました"
			if commands[0] == "iot" {
				if err := interactive.IotInvoke(commands[1], commands[2]); err != nil {
					responseText = fmt.Sprintf("処置に失敗しました: %v", err)
				}
			} else if commands[0] == "ec2" {
				if err := interactive.Ec2Invoke(commands[1], commands[2]); err != nil {
					responseText = fmt.Sprintf("処置に失敗しました: %v", err)
				}
			} else {
				responseText = "対象が未定義です"
			}
			api.PostMessage(channelId, slack.MsgOptionText(responseText, false))
		}
	})
}
