package main

import (
	"context"
	"net/url"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// Client implements the wwrclt.Client interface
type Client struct {
	clt wwrclt.Client
}

// NewClient creates a new autoconnected client instance
func NewClient(
	serverAddr url.URL,
	defaultReqTimeo time.Duration,
	transport wwr.ClientTransport,
) *Client {
	clt := &Client{}

	// Initialize client
	client, err := wwrclt.NewClient(
		serverAddr,
		clt,
		wwrclt.Options{
			// Default timeout for timed requests
			DefaultRequestTimeout: defaultReqTimeo,
			MessageBufferSize:     1024 * 16,
		},
		transport,
	)
	if err != nil {
		panic(err)
	}

	clt.clt = client

	return clt
}

// OnDisconnected implements the wwrclt.Implementation interface
func (cl *Client) OnDisconnected() {}

// OnSessionClosed implements the wwrclt.Implementation interface
func (cl *Client) OnSessionClosed() {}

// OnSessionCreated implements the wwrclt.Implementation interface
func (cl *Client) OnSessionCreated(_ *wwr.Session) {}

// OnSignal implements the wwrclt.Implementation interface
func (cl *Client) OnSignal(_ wwr.Message) {}

// Request sends a request to the server and blocks until the reply is received
func (cl *Client) Request(payload wwr.Payload) (wwr.Reply, error) {
	return cl.clt.Request(context.Background(), nil, payload)
}

// Close closes the client
func (cl *Client) Close() {
	cl.clt.Close()
}
