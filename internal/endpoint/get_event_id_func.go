package endpoint

import (
	"github.com/followedwind/slackbot/internal/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"net/http"
	"os"
	"strings"
)

func GetEventIdImpl(r *http.Request, cli *slack.Client, ev *slackevents.EventsAPIEvent) (eventId string) {
	switch innerData := ev.InnerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		//export SLACK_BOT_USERID="<@U****>"
		text := strings.ReplaceAll(innerData.Text, os.Getenv("SLACK_BOT_USERID"), "")
		text = strings.TrimSpace(strings.ReplaceAll(text, "\u00a0", "")) //nbsp
		util.InfoLog(text, 0)
		if text == PASS_HANDLER_ID {
			eventId = text
		}
	}
	return
}
