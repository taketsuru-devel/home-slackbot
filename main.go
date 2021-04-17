package main

import (
	"fmt"
	"github.com/followedwind/slackbot/internal/endpoint"
	"net/http"
)

func main() {

	http.Handle("/homeiot-to-slackbot", &endpoint.HomeIotEndpoint{})
	http.Handle("/events-endpoint", &endpoint.EventEndpoint{})
	http.Handle("/interactive", &endpoint.InteractiveEndpoint{})
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":13000", nil)

}
