package endpoint

import (
	"encoding/json"
	"fmt"
	"github.com/followedwind/slackbot/internal/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"io/ioutil"
	"net/http"
	"os"
)

type EventEndpoint struct{}

func (h *EventEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"), slack.OptionDebug(true))
	signingSecret := os.Getenv("SIGNING_SECRET")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := sv.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := sv.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}
	util.DebugLog(fmt.Sprintf("%#v\n", eventsAPIEvent.InnerEvent.Data))
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			//api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			api.PostMessage(ev.Channel, commandList())
		}
	}
}
