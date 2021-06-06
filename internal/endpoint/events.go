package endpoint

import (
	"encoding/json"
	"fmt"
	"github.com/followedwind/slackbot/internal/util"
	"github.com/slack-go/slack/slackevents"
	"net/http"
)

type EventEndpoint struct{}

func (h *EventEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := util.SlackRequestPreprocess(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(*body), slackevents.OptionNoVerifyToken())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(*body, &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	} else if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		util.DebugLog(fmt.Sprintf("%#v\n", innerEvent))
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			//api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			util.GetSlackClient().PostMessage(ev.Channel, commandList())
		}
	}
}
