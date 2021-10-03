package endpoint

import (
	"fmt"
	"net/http"

	"github.com/followedwind/slackbot/internal/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/taketsuru-devel/gorilla-microservice-skeleton/slackwrap"
)

type DefaultEventHandler struct{}

func (e *DefaultEventHandler) EventHandle() slackwrap.EventHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, cli *slack.Client, ev *slackevents.EventsAPIEvent) (interrupt bool, err error) {
		interrupt = true
		innerEvent := ev.InnerEvent
		util.DebugLog(fmt.Sprintf("%#v\n", innerEvent), 0)
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			_, _, err = cli.PostMessage(ev.Channel, commandList())
		default:
			err = fmt.Errorf("unsupported event")
		}
		return
	}
}
