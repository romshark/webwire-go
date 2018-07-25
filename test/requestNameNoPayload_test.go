package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestRequestNameNoPayload tests named requests without a payload
func TestRequestNameNoPayload(t *testing.T) {
	// Initialize server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ *webwire.Client,
				msg webwire.Message,
			) (webwire.Payload, error) {
				// Expect a named request
				msgName := msg.Name()
				if msgName != "n" {
					t.Errorf("Unexpected request name: %s", msgName)
				}

				// Expect no payload to arrive
				msgPayloadData := msg.Payload().Data()
				if len(msgPayloadData) > 0 {
					t.Errorf(
						"Unexpected received payload: %d",
						len(msgPayloadData),
					)
				}

				return nil, nil
			},
		},
		webwire.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	// Send a named binary request without a payload and await reply
	if _, err := client.connection.Request("n", webwire.NewPayload(
		webwire.EncodingBinary,
		nil,
	)); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Send a UTF16 encoded named binary request without a payload
	if _, err := client.connection.Request("n", webwire.NewPayload(
		webwire.EncodingUtf16,
		nil,
	)); err != nil {
		t.Fatalf("Request failed: %s", err)
	}
}
