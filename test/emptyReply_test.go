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
				_ wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Return empty reply
				return nil, nil
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
		callbackPoweredClientHooks{},
	)

	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send request and await reply
	reply, err := client.connection.Request(
		context.Background(),
		"",
		wwr.NewPayload(wwr.EncodingBinary, []byte("test")),
	)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify reply is empty
	replyEncoding := reply.Encoding()
	if replyEncoding != wwr.EncodingBinary {
		t.Fatalf(
			"Expected empty binary reply, but encoding was: %s",
			replyEncoding.String(),
		)
	}
	replyData := reply.Data()
	if len(replyData) > 0 {
		t.Fatalf("Expected empty binary reply, but payload was: %v", replyData)
	}
}
