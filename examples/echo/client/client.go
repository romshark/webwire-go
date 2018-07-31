package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	wwr "github.com/qbeon/webwire-go"

	wwrclt "github.com/qbeon/webwire-go/client"
)

// EchoClient implements the wwrclt.Implementation interface
type EchoClient struct {
	connection *wwrclt.Client
}

// NewEchoClient constructs and returns a new echo client instance
func NewEchoClient(serverAddr string) *EchoClient {
	newEchoClient := &EchoClient{}

	// Initialize connection
	newEchoClient.connection = wwrclt.NewClient(
		serverAddr,
		newEchoClient,
		wwrclt.Options{
			// Default timeout for timed requests
			DefaultRequestTimeout: 10 * time.Second,
			ReconnectionInterval:  2 * time.Second,
		},
	)

	return newEchoClient
}

// OnDisconnected implements the wwrclt.Implementation interface
func (clt *EchoClient) OnDisconnected() {}

// OnSessionClosed implements the wwrclt.Implementation interface
func (clt *EchoClient) OnSessionClosed() {}

// OnSessionCreated implements the wwrclt.Implementation interface
func (clt *EchoClient) OnSessionCreated(_ *wwr.Session) {}

// OnSignal implements the wwrclt.Implementation interface
func (clt *EchoClient) OnSignal(_ wwr.Payload) {}

// Request sends a message to the server and returns the reply.
// panics if the request fails for whatever reason
func (clt *EchoClient) Request(
	ctx context.Context,
	message string,
) wwr.Payload {
	// Define a payload to be sent to the server, use default binary encoding
	payload := wwr.NewPayload(wwr.EncodingBinary, []byte(message))

	log.Printf(
		"Sent request:   '%s' (%d)",
		string(payload.Data()),
		len(payload.Data()),
	)

	// Send request and await reply
	reply, err := clt.connection.Request(ctx, "", payload)
	if err != nil {
		panic(fmt.Errorf("Request failed: %s", err))
	}

	return reply
}

var serverAddr = flag.String("addr", ":8081", "server address")

func main() {
	// Parse command line arguments
	flag.Parse()

	// Initialize a new echo client instance
	echoClient := NewEchoClient(*serverAddr)

	// Send request and await reply
	reply := echoClient.Request(context.Background(), "hey, server!")

	log.Printf(
		"Received reply: '%s' (%d)",
		string(reply.Data()),
		len(reply.Data()),
	)
}
