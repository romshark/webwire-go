package test

import (
	"context"
	"os"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientRequestRegisterOnTimeout verifies the request register of the client
// is correctly updated when the request times out
func TestClientRequestRegisterOnTimeout(t *testing.T) {
	var client webwireClient.Client

	// Initialize webwire server given only the request
	server := setupServer(
		t,
		webwire.Hooks{
			OnRequest: func(ctx context.Context) ([]byte, *webwire.Error) {
				// Verify pending requests
				pendingReqs := client.Requests()
				if pendingReqs != 1 {
					t.Errorf("Unexpected pending requests: %d", pendingReqs)
					return nil, nil
				}

				// Wait until the request times out
				time.Sleep(300 * time.Millisecond)
				return nil, nil
			},
		},
	)
	go server.Run()

	// Initialize client
	client = webwireClient.NewClient(
		server.Addr,
		webwireClient.Hooks{},
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)

	// Connect the client to the server
	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect to the server: %s", err)
	}

	// Verify pending requests
	pendingReqsBeforeReq := client.Requests()
	if pendingReqsBeforeReq != 0 {
		t.Fatalf("Unexpected pending requests: %d", pendingReqsBeforeReq)
	}

	// Send request and await reply
	_, err := client.TimedRequest([]byte("t"), 200*time.Millisecond)
	if err == nil {
		t.Fatalf("Request must have failed (timeout)")
	}

	// Verify pending requests
	pendingReqsAfterReq := client.Requests()
	if pendingReqsAfterReq != 0 {
		t.Fatalf("Unexpected pending requests: %d", pendingReqsAfterReq)
	}
}
