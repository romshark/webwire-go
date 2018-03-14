package test

import (
	"context"
	"os"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientSignal verifies the server is connectable
// and receives signals correctly
func TestClientSignal(t *testing.T) {
	expectedSignalPayload := webwire.Payload{
		Encoding: webwire.EncodingUtf8,
		Data:     []byte("webwire_test_SIGNAL_payload"),
	}
	signalArrived := NewPending(1, 1*time.Second, true)

	// Initialize webwire server given only the signal handler
	_, addr := setupServer(
		t,
		webwire.ServerOptions{
			Hooks: webwire.Hooks{
				OnSignal: func(ctx context.Context) {
					// Extract signal message from the context
					msg := ctx.Value(webwire.Msg).(webwire.Message)

					// Verify signal payload
					comparePayload(
						t,
						"client signal",
						expectedSignalPayload,
						msg.Payload,
					)

					// Synchronize, notify signal arrival
					signalArrived.Done()
				},
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

	// Send signal
	err := client.Signal("", expectedSignalPayload)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Synchronize, await signal arrival
	if err := signalArrived.Wait(); err != nil {
		t.Fatal("Signal didn't arrive")
	}
}
