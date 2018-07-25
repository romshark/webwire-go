package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestEmptyReplyUtf16 verifies empty UTF16 encoded reply acceptance
func TestEmptyReplyUtf16(t *testing.T) {
	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ *wwr.Client,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Return empty reply
				return wwr.NewPayload(wwr.EncodingUtf16, nil), nil
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
	reply, err := client.connection.Request("", wwr.NewPayload(
		wwr.EncodingBinary,
		[]byte("test"),
	))
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify reply is empty
	replyEncoding := reply.Encoding()
	if replyEncoding != wwr.EncodingUtf16 {
		t.Fatalf(
			"Expected empty UTF16 reply, but encoding was: %s",
			replyEncoding.String(),
		)
	}
	replyData := reply.Data()
	if len(replyData) > 0 {
		t.Fatalf("Expected empty UTF16 reply, but payload was: %v", replyData)
	}
}
