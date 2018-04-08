package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientRequestRegisterOnTimeout verifies the request register
// of the client is correctly updated when the request times out
func TestClientRequestRegisterOnTimeout(t *testing.T) {
	var connection *webwireClient.Client

	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ *webwire.Client,
				_ *webwire.Message,
			) (webwire.Payload, error) {
				// Verify pending requests
				pendingReqs := connection.PendingRequests()
				if pendingReqs != 1 {
					t.Errorf("Unexpected pending requests: %d", pendingReqs)
					return webwire.Payload{}, nil
				}

				// Wait until the request times out
				time.Sleep(300 * time.Millisecond)
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
	_, reqErr := client.connection.TimedRequest(
		"",
		webwire.Payload{Data: []byte("t")},
		200*time.Millisecond,
	)
	if _, isTimeoutErr := reqErr.(webwire.ReqTimeoutErr); !isTimeoutErr {
		t.Fatalf("Request must have failed (timeout)")
	}

	// Verify pending requests
	pendingReqsAfterReq := client.connection.PendingRequests()
	if pendingReqsAfterReq != 0 {
		t.Fatalf("Unexpected pending requests: %d", pendingReqsAfterReq)
	}
}
