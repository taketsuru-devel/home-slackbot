package main

import (
	"flag"
	"fmt"
	"github.com/followedwind/slackbot/internal/endpoint"
	"github.com/followedwind/slackbot/internal/util"
	"github.com/taketsuru-devel/gorilla-microservice-skeleton/serverwrap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	prettyLog := flag.Bool("pretty_log", false, "sets pretty log")
	flag.Parse()
	util.InitLog(*debug, *prettyLog)
	util.InitSlackClient(false, nil, nil)

	server := serverwrap.NewServer(":13000")

	server.AddHandle("/homeiot-to-slackbot", &endpoint.HomeIotEndpoint{}).Methods("POST")
	server.AddHandle("/events-endpoint", &endpoint.EventEndpoint{}).Methods("POST")
	server.AddHandle("/interactive", &endpoint.InteractiveEndpoint{}).Methods("POST")

	server.Start()
	defer server.Stop(60)

	waitSignal()
}

func waitSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	util.DebugLog(fmt.Sprintf("terminate signal(%d) received", <-quit), 0)
	close(quit)
}
