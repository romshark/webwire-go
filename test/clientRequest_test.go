package test

import (
	"context"
	"os"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientRequest verifies the server is connectable,
// receives requests and answers them correctly
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
	_, addr := setupServer(
		t,
		webwire.Hooks{
			OnRequest: func(ctx context.Context) (webwire.Payload, error) {
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
	)

	// Initialize client
	client := webwireClient.NewClient(
		addr,
		webwireClient.Hooks{},
		5*time.Second,
		os.Stdout,
		os.Stderr,
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
