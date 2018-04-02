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
			onRequest: func(_ context.Context) (wwr.Payload, error) {
				// Return empty reply
				return wwr.Payload{
					Encoding: wwr.EncodingUtf16,
				}, nil
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := webwireClient.NewClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)

	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send request and await reply
	reply, err := client.Request("", wwr.Payload{Data: []byte("test")})
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify reply is empty
	if reply.Encoding != wwr.EncodingUtf16 {
		t.Fatalf(
			"Expected empty UTF16 reply, but encoding was: %s",
			reply.Encoding.String(),
		)
	}
	if len(reply.Data) > 0 {
		t.Fatalf("Expected empty UTF16 reply, but payload was: %v", reply.Data)
	}
}
