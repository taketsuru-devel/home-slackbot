package endpoint

import (
	"encoding/json"
	"fmt"
	"github.com/followedwind/slackbot/internal/interactive"
	"github.com/followedwind/slackbot/internal/util"
	"github.com/slack-go/slack"
	"net/http"
	"strings"
)

type InteractiveEndpoint struct{}

func (i *InteractiveEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//情報取得
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		util.ErrorLog(fmt.Sprintf("Could not parse action response JSON: %v", err))
		//ここでエラーだとChannelの取得もできない
		return
	}
	channelId := payload.Channel.GroupConversation.Conversation.ID
	command := payload.ActionCallback.BlockActions[0].Value
	commandDisp := payload.ActionCallback.BlockActions[0].Text.Text

	util.DebugLog(fmt.Sprintf("command:%#v", payload.ActionCallback.BlockActions[0].Value))

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
