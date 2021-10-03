package endpoint

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/followedwind/slackbot/internal/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/taketsuru-devel/gorilla-microservice-skeleton/slackwrap"
	"net/http"
	"strings"
	"time"
)

const PASS_HANDLER_ID = "pw"

type PassHandler struct{}

func (p *PassHandler) GetEventHandler() slackwrap.EventHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, cli *slack.Client, ev *slackevents.EventsAPIEvent) (interrupt bool, err error) {
		interrupt = true
		switch innerData := ev.InnerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			//一覧表示してlistにしてイベントを作る
			//Dynamo getAll and show list
			svc := dynamodb.New(session.New(), &aws.Config{Region: aws.String("ap-northeast-1")})
			input := &dynamodb.ScanInput{
				AttributesToGet: []*string{aws.String("Service")},
				TableName:       aws.String("PassDb"),
			}
			if result, dynamoErr := svc.Scan(input); dynamoErr != nil {
				err = dynamoErr
			} else {
				//一覧表示してlistにしてイベントを作る
				selectDefault := slack.NewTextBlockObject("plain_text", "未選択", false, false)
				opts := make([]*slack.OptionBlockObject, 0, len(result.Items))
				for _, data := range result.Items {
					txt := *data["Service"].S
					txtObj := slack.NewTextBlockObject("plain_text", txt, false, false)
					opts = append(opts, slack.NewOptionBlockObject(txt, txtObj, txtObj))
				}
				opte := slack.NewOptionsSelectBlockElement("static_select", selectDefault, PASS_HANDLER_ID, opts...)
				notice := slack.NewTextBlockObject("plain_text", "以下から選択してください", false, false)
				mbk := slack.NewSectionBlock(notice, nil, slack.NewAccessory(opte))
				if _, ts, postErr := cli.PostMessage(innerData.Channel, slack.MsgOptionBlocks(mbk)); postErr != nil {
					err = postErr
				} else {
					util.InfoLog(fmt.Sprintf("selection ts:%s", ts), 0)
				}
			}
		default:
			err = fmt.Errorf("unsupported type: %v", innerData)
		}

		return
	}
}

func (p *PassHandler) GetEventId() string {
	return PASS_HANDLER_ID
}

func (p *PassHandler) GetBlockActionHandler() slackwrap.BlockActionHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, cli *slack.Client, ic *slack.InteractionCallback, b *slack.BlockAction) (err error) {
		channelId := ic.Channel.GroupConversation.Conversation.ID
		cli.PostMessage(channelId, slack.MsgOptionText(fmt.Sprintf("%sを受け付けました", b.SelectedOption.Value), false))
		//Dynamo getAll and show list
		svc := dynamodb.New(session.New(), &aws.Config{Region: aws.String("ap-northeast-1")})
		input := &dynamodb.QueryInput{
			TableName:                 aws.String("PassDb"),
			ProjectionExpression:      aws.String("PassData"),
			KeyConditionExpression:    aws.String("Service = :v1"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":v1": &dynamodb.AttributeValue{S: aws.String(b.SelectedOption.Value)}},
		}
		if result, dynamoErr := svc.Query(input); dynamoErr != nil {
			err = dynamoErr
		} else {
			//一覧表示してlistにしてイベントを作る
			for _, data := range result.Items {
				dataMap := data["PassData"].M
				strbuf := make([]string, 0, len(dataMap))
				for k, v := range dataMap {
					strbuf = append(strbuf, fmt.Sprintf("%s: %s", k, *v.S))
				}
				if _, ts, postErr := cli.PostMessage(channelId, slack.MsgOptionText(strings.Join(strbuf, "\n"), false)); err != nil {
					err = postErr
				} else {
					//連続だとエラーが出るので少し寝る
					time.Sleep(100 * time.Millisecond)
					if _, _, postErr := cli.PostMessage(channelId, slack.MsgOptionText("60秒後に削除されます", false), slack.MsgOptionTS(ts)); postErr != nil {
						err = postErr
					}
					go func() {
						time.Sleep(60 * time.Second)
						cli.DeleteMessage(channelId, ts)
					}()
				}
			}
		}
		return
	}
}
