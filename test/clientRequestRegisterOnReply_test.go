package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientRequestRegisterOnReply verifies the request register of the client
// is correctly updated when the request is successfully fulfilled
func TestClientRequestRegisterOnReply(t *testing.T) {
	var connection *webwireClient.Client

	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ webwire.Connection,
				_ webwire.Message,
			) (webwire.Payload, error) {
				// Verify pending requests
				pendingReqs := connection.PendingRequests()
				if pendingReqs != 1 {
					t.Errorf("Unexpected pending requests: %d", pendingReqs)
					return nil, nil
				}

				// Wait until the request times out
				time.Sleep(300 * time.Millisecond)
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
	connection = client.connection

	// Connect the client to the server
	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect to the server: %s", err)
	}

	// Verify pending requests
	pendingReqsBeforeReq := client.connection.PendingRequests()
	if pendingReqsBeforeReq != 0 {
		t.Fatalf("Unexpected pending requests: %d", pendingReqsBeforeReq)
	}

	// Send request and await reply
	_, err := client.connection.Request("", webwire.NewPayload(
		webwire.EncodingBinary,
		[]byte("t"),
	))
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify pending requests
	pendingReqsAfterReq := client.connection.PendingRequests()
	if pendingReqsAfterReq != 0 {
		t.Fatalf("Unexpected pending requests: %d", pendingReqsAfterReq)
	}
}
