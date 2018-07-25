package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestRequestNoNameOnlyPayload tests requests without a name but only a payload
func TestRequestNoNameOnlyPayload(t *testing.T) {
	expectedRequestPayload := webwire.NewPayload(
		webwire.EncodingUtf8,
		[]byte("3"),
	)
	expectedRequestPayloadUtf16 := webwire.NewPayload(
		webwire.EncodingUtf16,
		[]byte("12"),
	)

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
				if len(msgName) > 0 {
					t.Errorf("Unexpected request name: %s", msgName)
				}

				msgPayload := msg.Payload()
				if msgPayload.Encoding() == webwire.EncodingUtf16 {
					// Verify request payload
					comparePayload(
						t,
						"client request (UTF16)",
						expectedRequestPayloadUtf16,
						msgPayload,
					)
				} else {
					// Verify request payload
					comparePayload(
						t,
						"client request",
						expectedRequestPayload,
						msgPayload,
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

	// Send an unnamed binary request with a payload and await reply
	_, err := client.connection.Request("", expectedRequestPayload)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Send an unnamed UTF16 encoded binary request with a payload
	_, err = client.connection.Request("", expectedRequestPayloadUtf16)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}
}
