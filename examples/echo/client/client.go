package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/qbeon/webwire-go"

	webwireClient "github.com/qbeon/webwire-go/client"
)

var serverAddr = flag.String("addr", ":8081", "server address")

func main() {
	// Parse command line arguments
	flag.Parse()

	// Define a payload to be sent to the server, use default binary encoding
	payload := webwire.Payload{
		Data: []byte("hey server!"),
	}

	// Initialize client
	client := webwireClient.NewClient(
		*serverAddr,
		webwireClient.Options{
			// No hooks required in this example
			Hooks: webwireClient.Hooks{},
			// Default timeout for timed requests
			DefaultRequestTimeout: 5 * time.Second,
		},
	)

	log.Printf("Connect to %s", *serverAddr)

	if err := client.Connect(); err != nil {
		panic(fmt.Errorf("Couldn't connect to the server: %s", err))
	}

	log.Printf("Send request: '%s' (%d)", string(payload.Data), len(payload.Data))

	// Send request and await reply
	reply, err := client.Request("", payload)
	if err != nil {
		panic(fmt.Errorf("Request failed: %s", err))
	}

	log.Printf("Received reply: '%s' (%d)", string(reply.Data), len(reply.Data))
}
