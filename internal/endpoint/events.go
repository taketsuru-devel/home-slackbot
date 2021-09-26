package endpoint

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/followedwind/slackbot/internal/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/taketsuru-devel/gorilla-microservice-skeleton/slackwrap"
)

func GetEventHandler() *slackwrap.EventSubscribeEndpoint {
	signingSecret := os.Getenv("SIGNING_SECRET")
	return &slackwrap.EventSubscribeEndpoint{
		Handler:       &eventHandler{},
		SigningSecret: &signingSecret,
	}
}

type eventHandler struct{}

func (h *eventHandler) Handle() slackwrap.EventSubscribeHandlerFunc {
	return (func(w http.ResponseWriter, r *http.Request, ev *slackevents.EventsAPIEvent) {
		util.InfoLog(fmt.Sprintf("%v", ev.InnerEvent.Data), 0)

		innerEvent := ev.InnerEvent
		util.DebugLog(fmt.Sprintf("%#v\n", innerEvent), 0)
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			text := strings.ReplaceAll(ev.Text, "<@U01MY2T85LM>", "")
			text = strings.TrimSpace(strings.ReplaceAll(text, "\u00a0", "")) //nbsp
			if text == "" {
				util.GetSlackClient().PostMessage(ev.Channel, commandList())
			} else if text == "pw" {
				//Dynamo getAll and show list
				svc := dynamodb.New(session.New(), &aws.Config{Region: aws.String("ap-northeast-1")})
				input := &dynamodb.ScanInput{
					AttributesToGet: []*string{aws.String("Service")},
					TableName:       aws.String("PassDb"),
				}
				if result, err := svc.Scan(input); err != nil {
					fmt.Println(err)
				} else {
					//一覧表示してlistにしてイベントを作る
					selectDefault := slack.NewTextBlockObject("plain_text", "未選択", false, false)
					opts := make([]*slack.OptionBlockObject, 0, len(result.Items))
					for _, data := range result.Items {
						txt := *data["Service"].S
						txtObj := slack.NewTextBlockObject("plain_text", txt, false, false)
						opts = append(opts, slack.NewOptionBlockObject(txt, txtObj, txtObj))
					}
					opte := slack.NewOptionsSelectBlockElement("static_select", selectDefault, "pwselect", opts...)
					notice := slack.NewTextBlockObject("plain_text", "以下から選択してください", false, false)
					mbk := slack.NewSectionBlock(notice, nil, slack.NewAccessory(opte))
					if _, ts, err := util.GetSlackClient().PostMessage(ev.Channel, slack.MsgOptionBlocks(mbk)); err != nil {
						fmt.Println(err)
					} else {
						util.InfoLog(fmt.Sprintf("selection ts:%s", ts), 0)
					}
				}
			}
		}
	})
}
