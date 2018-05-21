package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestRequestEmpty tests empty requests without a name and without a payload
func TestRequestEmpty(t *testing.T) {
	// Initialize server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ *webwire.Client,
				msg *webwire.Message,
			) (webwire.Payload, error) {
				// Expect the following request to not even arrive
				t.Error("Not expected but reached")
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

	// Send request without a name and without a payload.
	// Expect a protocol error in return not sending the invalid request off
	_, err := client.connection.Request("", webwire.Payload{})
	if _, isProtoErr := err.(webwire.ProtocolErr); !isProtoErr {
		t.Fatalf("Expected a protocol error, got: %v", err)
	}
}
