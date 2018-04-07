package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientSignal tests client-side signals with UTF8 encoded payloads
func TestClientSignal(t *testing.T) {
	expectedSignalPayload := webwire.Payload{
		Encoding: webwire.EncodingUtf8,
		Data:     []byte("webwire_test_SIGNAL_payload"),
	}
	signalArrived := newPending(1, 1*time.Second, true)

	// Initialize webwire server given only the signal handler
	server := setupServer(
		t,
		&serverImpl{
			onSignal: func(
				_ context.Context,
				_ *webwire.Client,
				msg *webwire.Message,
			) {

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
		webwire.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, nil, nil, nil,
	)

	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send signal
	err := client.connection.Signal("", expectedSignalPayload)
	if err != nil {
		t.Fatalf("Couldn't send signal: %s", err)
	}

	// Synchronize, await signal arrival
	if err := signalArrived.Wait(); err != nil {
		t.Fatal("Signal wasn't processed")
	}
}
