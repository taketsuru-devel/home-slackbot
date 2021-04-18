package util

import (
	"github.com/slack-go/slack"
	"io/ioutil"
	"net/http"
	"os"
)

//https://pkg.go.dev/github.com/slack-go/slack

var clientOptions []slack.Option

func GetSlackClient() *slack.Client {
	return slack.New(os.Getenv("SLACK_BOT_TOKEN"), clientOptions...)
}

func VerifySlackSecret(r *http.Request) error {
	signingSecret := os.Getenv("SIGNING_SECRET")
	sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if _, err := sv.Write(body); err != nil {
		return err
	}
	if err := sv.Ensure(); err != nil {
		return err
	}
	return nil
}

func InitSlackClient(debug bool, test bool) {
	options := make([]slack.Option, 0, 2)
	if debug {
		options = append(options, slack.OptionDebug(true))
	}
	if test {
		//tekito
		options = append(options, slack.OptionAPIURL("http://localhost:6000"))
	}
	clientOptions = options
}
