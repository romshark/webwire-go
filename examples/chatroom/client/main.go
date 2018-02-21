package main

import (
	"os"
	"fmt"
	"log"
	"flag"
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
		// Address of the webwire server
		*serverAddr,
		onSignal,
		onSessionCreated,
		nil,

		// Default timeout for timed requests
		5 * time.Second,
		
		// Log writers
		os.Stdout,
		os.Stderr,
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
