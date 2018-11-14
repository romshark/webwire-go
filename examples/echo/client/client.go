package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	wwrfasthttp "github.com/qbeon/webwire-go/transport/fasthttp"
)

// EchoClient implements the wwrclt.Implementation interface
type EchoClient struct {
	connection wwrclt.Client
}

// NewEchoClient constructs and returns a new echo client instance
func NewEchoClient(serverAddr url.URL) (*EchoClient, error) {
	newEchoClient := &EchoClient{}

	// Initialize connection
	connection, err := wwrclt.NewClient(
		serverAddr,
		newEchoClient,
		wwrclt.Options{
			// Default timeout for timed requests
			DefaultRequestTimeout: 10 * time.Second,
			ReconnectionInterval:  2 * time.Second,
		},
		&wwrfasthttp.ClientTransport{},
	)
	if err != nil {
		return nil, err
	}

	newEchoClient.connection = connection

	return newEchoClient, nil
}

// OnDisconnected implements the wwrclt.Implementation interface
func (clt *EchoClient) OnDisconnected() {}

// OnSessionClosed implements the wwrclt.Implementation interface
func (clt *EchoClient) OnSessionClosed() {}

// OnSessionCreated implements the wwrclt.Implementation interface
func (clt *EchoClient) OnSessionCreated(_ *wwr.Session) {}

// OnSignal implements the wwrclt.Implementation interface
func (clt *EchoClient) OnSignal(_ wwr.Message) {}

// Request sends a message to the server and returns the reply.
// panics if the request fails for whatever reason
func (clt *EchoClient) Request(message string) []byte {
	// Define a payload to be sent to the server, use default binary encoding
	payload := wwr.Payload{Data: []byte(message)}

	log.Printf(
		"Sent request:   '%s' (%d)",
		string(payload.Data),
		len(payload.Data),
	)

	// Send request and await reply
	reply, err := clt.connection.Request(context.Background(), nil, payload)
	if err != nil {
		panic(fmt.Errorf("Request failed: %s", err))
	}

	// Copy the reply payload
	pld := reply.Payload()
	data := make([]byte, len(pld))
	copy(data, pld)

	// Close the reply to release the buffer
	reply.Close()

	return data
}

var serverAddr = flag.String("addr", ":8081", "server address")

func main() {
	// Parse command line arguments
	flag.Parse()

	// Initialize a new echo client instance
	echoClient, err := NewEchoClient(url.URL{Host: *serverAddr})
	if err != nil {
		panic(err)
	}

	// Send request and await reply
	reply := echoClient.Request("hey, server!")

	log.Printf(
		"Received reply: '%s' (%d)",
		string(reply),
		len(reply),
	)
}
