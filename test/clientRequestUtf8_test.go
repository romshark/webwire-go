package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientRequest tests requests with UTF8 encoded payloads
func TestClientRequest(t *testing.T) {
	expectedRequestPayload := webwire.Payload{
		Encoding: webwire.EncodingUtf8,
		Data:     []byte("webwire_test_REQUEST_payload"),
	}
	expectedReplyPayload := webwire.Payload{
		Encoding: webwire.EncodingUtf8,
		Data:     []byte("webwire_test_RESPONSE_message"),
	}

	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(ctx context.Context) (webwire.Payload, error) {
				// Extract request message from the context
				msg := ctx.Value(webwire.Msg).(webwire.Message)

				// Verify request payload
				comparePayload(
					t,
					"client request",
					expectedRequestPayload,
					msg.Payload,
				)
				return expectedReplyPayload, nil
			},
		},
		webwire.ServerOptions{},
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
	reply, err := client.Request("", expectedRequestPayload)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify reply
	comparePayload(t, "server reply", expectedReplyPayload, reply)
}
