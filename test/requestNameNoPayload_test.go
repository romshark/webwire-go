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
				msg *webwire.Message,
			) (webwire.Payload, error) {
				// Expect a named request
				if msg.Name != "n" {
					t.Errorf("Unexpected request name: %s", msg.Name)
				}

				// Expect no payload to arrive
				if len(msg.Payload.Data) > 0 {
					t.Errorf(
						"Unexpected received payload: %d",
						len(msg.Payload.Data),
					)
				}

				return webwire.Payload{}, nil
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
	if _, err := client.connection.Request("n", webwire.Payload{
		Encoding: webwire.EncodingBinary,
		Data:     nil,
	}); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Send a UTF16 encoded named binary request without a payload
	if _, err := client.connection.Request("n", webwire.Payload{
		Encoding: webwire.EncodingUtf16,
		Data:     nil,
	}); err != nil {
		t.Fatalf("Request failed: %s", err)
	}
}
