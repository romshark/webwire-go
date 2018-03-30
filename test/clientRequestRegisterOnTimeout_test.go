package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientRequestRegisterOnTimeout verifies the request register of the client
// is correctly updated when the request times out
func TestClientRequestRegisterOnTimeout(t *testing.T) {
	var client *webwireClient.Client

	// Initialize webwire server given only the request
	server := setupServer(
		t,
		webwire.ServerOptions{
			Hooks: webwire.Hooks{
				OnRequest: func(ctx context.Context) (webwire.Payload, error) {
					// Verify pending requests
					pendingReqs := client.PendingRequests()
					if pendingReqs != 1 {
						t.Errorf("Unexpected pending requests: %d", pendingReqs)
						return webwire.Payload{}, nil
					}

					// Wait until the request times out
					time.Sleep(300 * time.Millisecond)
					return webwire.Payload{}, nil
				},
			},
		},
	)

	// Initialize client
	client = webwireClient.NewClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)

	// Connect the client to the server
	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect to the server: %s", err)
	}

	// Verify pending requests
	pendingReqsBeforeReq := client.PendingRequests()
	if pendingReqsBeforeReq != 0 {
		t.Fatalf("Unexpected pending requests: %d", pendingReqsBeforeReq)
	}

	// Send request and await reply
	_, reqErr := client.TimedRequest(
		"",
		webwire.Payload{Data: []byte("t")},
		200*time.Millisecond,
	)
	if _, isTimeoutErr := reqErr.(webwire.ReqTimeoutErr); !isTimeoutErr {
		t.Fatalf("Request must have failed (timeout)")
	}

	// Verify pending requests
	pendingReqsAfterReq := client.PendingRequests()
	if pendingReqsAfterReq != 0 {
		t.Fatalf("Unexpected pending requests: %d", pendingReqsAfterReq)
	}
}
