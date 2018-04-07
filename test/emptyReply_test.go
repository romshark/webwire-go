package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestEmptyReply verifies empty binary reply acceptance
func TestEmptyReply(t *testing.T) {
	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ *wwr.Client,
				_ *wwr.Message,
			) (wwr.Payload, error) {
				// Return empty reply
				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, nil, nil, nil,
	)

	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send request and await reply
	reply, err := client.connection.Request("", wwr.Payload{
		Data: []byte("test"),
	})
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify reply is empty
	if reply.Encoding != wwr.EncodingBinary {
		t.Fatalf(
			"Expected empty binary reply, but encoding was: %s",
			reply.Encoding.String(),
		)
	}
	if len(reply.Data) > 0 {
		t.Fatalf("Expected empty binary reply, but payload was: %v", reply.Data)
	}
}
