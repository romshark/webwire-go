package main

import (
	"flag"
	"time"

	webwireClient "github.com/qbeon/webwire-go/client"
)

var client *webwireClient.Client

var serverAddr = flag.String("addr", ":9090", "server address")
var password = flag.String("pass", "", "password")
var username = flag.String("name", "", "username")

func main() {
	// Parse command line arguments
	flag.Parse()

	// Initialize client
	client = webwireClient.NewClient(
		*serverAddr,
		webwireClient.Options{
			// Address of the webwire server
			Hooks: webwireClient.Hooks{
				OnServerSignal:   onSignal,
				OnSessionCreated: onSessionCreated,
			},
			// Default timeout for timed requests
			DefaultRequestTimeout: 5 * time.Second,
		},
	)
	defer client.Close()

	// Authenticate
	if *username != "" && *password != "" {
		authenticate()
	}

	// Main input loop
	mainLoop()
}
