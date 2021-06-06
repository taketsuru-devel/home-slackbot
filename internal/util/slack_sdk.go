package util

import (
	"github.com/slack-go/slack"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//https://pkg.go.dev/github.com/slack-go/slack

var clientOptions []slack.Option

func GetSlackClient() *slack.Client {
	return slack.New(os.Getenv("SLACK_BOT_TOKEN"), clientOptions...)
}

func SlackRequestPreprocess(r *http.Request) (*[]byte, error) {
	signingSecret := os.Getenv("SIGNING_SECRET")
	sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if _, err := sv.Write(body); err != nil {
		return nil, err
	}
	if err := sv.Ensure(); err != nil {
		return nil, err
	}
	return &body, nil
}

func InitSlackClient(debug bool, serverUrl *string, loggerWriter io.Writer) {
	options := make([]slack.Option, 0, 3)
	if debug {
		options = append(options, slack.OptionDebug(true))
	}
	if serverUrl != nil {
		options = append(options, slack.OptionAPIURL(*serverUrl))
	}
	if loggerWriter != nil {
		options = append(options, slack.OptionLog(log.New(loggerWriter, "slacktest", log.LstdFlags|log.Lshortfile)))
	}
	clientOptions = options
}
