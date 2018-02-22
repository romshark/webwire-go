package test

import (
	"context"
	"os"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwire_client "github.com/qbeon/webwire-go/client"
)

// TestClientSignal verifies the server is connectable
// and receives signals correctly
func TestClientSignal(t *testing.T) {
	expectedSignalPayload := []byte("webwire_test_SIGNAL_payload")
	wait := make(chan bool)

	// Initialize webwire server given only the signal handler
	server := setupServer(
		t,
		webwire.Hooks{
			OnSignal: func(ctx context.Context) {
				// Extract signal message from the context
				msg := ctx.Value(webwire.MESSAGE).(webwire.Message)

				// Verify signal payload
				comparePayload(
					t,
					"client signal",
					expectedSignalPayload,
					msg.Payload,
				)

				// Synchronize, notify signal arrival
				wait <- true
			},
		},
	)
	go server.Run()

	// Initialize client
	client := webwire_client.NewClient(
		server.Addr,
		nil, nil, nil,
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)

	// Send signal
	err := client.Signal(expectedSignalPayload)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Synchronize, await signal arrival
	<-wait
}
