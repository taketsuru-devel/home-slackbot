package util

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
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

func SlackRequestPreprocess(w http.ResponseWriter, r *http.Request) (*slackevents.EventsAPIEvent, error) {
	//bodyの取得
	/*
		//GetBody()はserver側では使えないみたい
		bodyReader, err := r.GetBody()
		if err != nil {
			return nil, err
		}
	*/
	bodyReader := r.Body
	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}
	defer bodyReader.Close()

	//署名の検証
	if err := verifySlackSecret(r.Header, &body); err != nil {
		return nil, err
	}

	//event型の取得
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	//slackからのAPI検証リクエストならここで処理する
	if done, err := handleAPIVerificationRequest(&eventsAPIEvent, &body, w); done {
		//API検証リクエストだったので済み
		return nil, err
	} else {
		//API検証リクエストじゃなかったので丸投げ
		return &eventsAPIEvent, err
	}
}

func verifySlackSecret(header http.Header, body *[]byte) error {
	signingSecret := os.Getenv("SIGNING_SECRET")
	sv, err := slack.NewSecretsVerifier(header, signingSecret)
	if err != nil {
		return err
	}
	if _, err := sv.Write(*body); err != nil {
		return err
	}
	if err := sv.Ensure(); err != nil {
		return err
	}
	return nil
}

func handleAPIVerificationRequest(ev *slackevents.EventsAPIEvent, body *[]byte, w http.ResponseWriter) (bool, error) {
	//slackからのAPI検証
	if ev.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(*body, &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return false, err
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
		return true, nil
	}
	return false, nil
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
