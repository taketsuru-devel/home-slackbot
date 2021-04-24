package endpoint

import (
	"fmt"
	"github.com/followedwind/slackbot/internal/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestInteractiveEndpoint(t *testing.T) {
	s := slacktest.NewTestServer()
	serverUrl := s.GetAPIURL()
	util.InitSlackClient(false, &serverUrl)
	s.Start()
	defer s.Stop()

	values := url.Values{}
	payload := slack.InteractionCallback{
		Channel: slack.Channel{
			GroupConversation: slack.GroupConversation{
				Conversation: slack.Conversation{
					ID: "test",
				},
			},
		},
		ActionCallback: slack.ActionCallbacks{
			BlockActions: []*slack.BlockAction{
				&slack.BlockAction{
					Value: "value:value",
					Text: slack.TextBlockObject{
						Text: "text",
					},
				},
			},
		},
	}
	payloadJson, _ := payload.MarshalJSON()
	values.Set("payload", string(payloadJson))
	req := httptest.NewRequest(http.MethodPost, "http://dummy", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()

	target := InteractiveEndpoint{}
	target.ServeHTTP(res, req)

	fmt.Println(req)
	fmt.Println(res)
	t.Errorf("test")
}
