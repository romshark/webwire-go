package test

import (
	"testing"
	"os"
	"time"
	"context"

	webwire "github.com/qbeon/webwire-go"
	webwire_client "github.com/qbeon/webwire-go/client"
)

// TestClientRequest verifies the server is connectable,
// receives requests and answers them correctly
func TestClientRequest(t *testing.T) {
	expectedRequestPayload := []byte("webwire_test_REQUEST_payload")
	expectedReplyPayload := []byte("webwire_test_RESPONSE_message")

	// Initialize webwire server given only the request
	server := setupServer(
		t,
		nil, nil,
		func(ctx context.Context) ([]byte, *webwire.Error) {
			// Extract request message from the context
			msg := ctx.Value(webwire.MESSAGE).(webwire.Message)

			// Verify request payload
			comparePayload(t, "client request", expectedRequestPayload, msg.Payload)
			return expectedReplyPayload, nil
		},
		nil, nil, nil, nil,
	)
	go server.Run()

	// Initialize client
	client := webwire_client.NewClient(
		server.Addr,
		nil,
		nil,
		5 * time.Second,
		os.Stdout,
		os.Stderr,
	)

	// Send request and await reply
	reply, err := client.Request(expectedRequestPayload)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify reply
	comparePayload(t, "server reply", expectedReplyPayload, reply)
}
