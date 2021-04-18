package main

import (
	"flag"
	"fmt"
	"github.com/followedwind/slackbot/internal/endpoint"
	"github.com/followedwind/slackbot/internal/serverwrap"
	"github.com/followedwind/slackbot/internal/util"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	prettyLog := flag.Bool("pretty_log", false, "sets pretty log")
	flag.Parse()
	util.InitLog(*debug, *prettyLog)

	server := serverwrap.NewServer(":13000")

	server.AddHandle("/homeiot-to-slackbot", &endpoint.HomeIotEndpoint{})
	server.AddHandle("/events-endpoint", &endpoint.EventEndpoint{})
	server.AddHandle("/interactive", &endpoint.InteractiveEndpoint{})

	server.Start()
	defer server.Stop()

	waitSignal()
}

func waitSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	util.DebugLog(fmt.Sprintf("terminate signal(%d) received", <-quit))
	close(quit)
}
