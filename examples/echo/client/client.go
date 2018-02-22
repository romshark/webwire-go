package main

import (
	"os"
	"fmt"
	"log"
	"flag"
	"time"

	webwireClient "github.com/qbeon/webwire-go/client"
)

var serverAddr = flag.String("addr", ":8081", "server address")

func main() {
	// Parse command line arguments
	flag.Parse()

	payload := []byte("hey server!")

	// Initialize client
	client := webwireClient.NewClient(
		// Address of the webwire server
		*serverAddr,

		// No hooks required in this example
		nil, nil, nil,

		// Default timeout for timed requests
		5 * time.Second,
		
		// Log writers
		os.Stdout,
		os.Stderr,
	)

	log.Printf("Connect to %s", *serverAddr)

	if err := client.Connect(); err != nil {
		panic(fmt.Errorf("Couldn't connect to the server: %s", err))
	}

	log.Printf("Send request: '%s' (%d)", string(payload), len(payload))

	// Send request and await reply
	reply, err := client.Request(payload)
	if err != nil {
		panic(fmt.Errorf("Request failed: %s", err))
	}

	log.Printf("Received reply: '%s' (%d)", string(reply), len(reply))
}
