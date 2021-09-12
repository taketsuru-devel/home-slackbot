package endpoint

import (
	"fmt"
	"github.com/followedwind/slackbot/internal/interactive"
	"github.com/followedwind/slackbot/internal/util"
	"github.com/slack-go/slack"
	"net/http"
	"strings"

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

		util.DebugLog(fmt.Sprintf("command:%#v", ic.ActionCallback.BlockActions[0].Value), 0)

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
	})
}
