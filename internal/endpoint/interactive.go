package endpoint

import (
	"encoding/json"
	"fmt"
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
	if err := util.IotInvoke(commands[0], commands[1]); err != nil {
		responseText = fmt.Sprintf("処置に失敗しました: %v", err)
	}
	api.PostMessage(channelId, slack.MsgOptionText(responseText, false))
}
