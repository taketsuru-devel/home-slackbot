package main

import (
	"fmt"
	"github.com/followedwind/slackbot/internal/endpoint"
	"github.com/followedwind/slackbot/internal/serverwrap"
	"os"
	"os/signal"
	"syscall"
)

func main() {

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
	fmt.Printf("terminate signal(%d) received\n", <-quit)
	close(quit)
}
