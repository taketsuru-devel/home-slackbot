package main

import (
	"os"

	"github.com/followedwind/slackbot/internal/endpoint"
	"github.com/followedwind/slackbot/internal/util"
	"github.com/slack-go/slack"
	"github.com/taketsuru-devel/gorilla-microservice-skeleton/serverwrap"
	"github.com/taketsuru-devel/gorilla-microservice-skeleton/skeletonutil"
	"github.com/taketsuru-devel/gorilla-microservice-skeleton/slackwrap"
)

func main() {
	util.InitLog(os.Getenv("LOG_DEBUG") == "True")
	util.InitSlackClient(false, nil, nil)

	server := serverwrap.NewServer(":13000")
	cli := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	slackSecret := os.Getenv("SIGNING_SECRET")
	f := slackwrap.NewSlackHandlerFactory(cli, &slackSecret, &endpoint.DefaultEventHandler{}, &endpoint.DefaultInteractiveHandler{})
	f.InitBlockAction(endpoint.GetEventIdImpl)
	f.RegisterBlockAction(&endpoint.PassHandler{})

	server.AddHandle("/homeiot-to-slackbot", &endpoint.HomeIotEndpoint{}).Methods("POST")
	server.AddHandle("/events-endpoint", f.CreateEventEndpoint()).Methods("POST")
	server.AddHandle("/interactive", f.CreateInteractiveEndpoint()).Methods("POST")

	server.Start()
	defer server.Stop(60)

	skeletonutil.WaitSignal()
}
