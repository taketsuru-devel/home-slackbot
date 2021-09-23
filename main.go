package main

import (
	"github.com/followedwind/slackbot/internal/endpoint"
	"github.com/followedwind/slackbot/internal/util"
	"github.com/taketsuru-devel/gorilla-microservice-skeleton/serverwrap"
	"github.com/taketsuru-devel/gorilla-microservice-skeleton/skeletonutil"
	"os"
)

func main() {
	util.InitLog(os.Getenv("LOG_DEBUG") == "True")
	util.InitSlackClient(false, nil, nil)

	server := serverwrap.NewServer(":13000")

	server.AddHandle("/homeiot-to-slackbot", &endpoint.HomeIotEndpoint{}).Methods("POST")
	server.AddHandle("/events-endpoint", endpoint.GetEventHandler()).Methods("POST")
	server.AddHandle("/interactive", endpoint.GetInteractiveHandler()).Methods("POST")

	server.Start()
	defer server.Stop(60)

	skeletonutil.WaitSignal()
}
