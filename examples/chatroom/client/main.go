package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	webwireClient "github.com/qbeon/webwire-go/client"
)

var client webwireClient.Client

var serverAddr = flag.String("addr", ":8081", "server address")
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

	// Connect to the server
	if err := client.Connect(); err != nil {
		panic(fmt.Errorf("Couldn't connect to the server: %s", err))
	}
	log.Printf("Connected to server %s", *serverAddr)

	// Authenticate
	if *username != "" && *password != "" {
		authenticate()
	}

	// Main input loop
	mainLoop()
}
