package test

import (
	"context"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientSignal tests client-side signals with UTF8 encoded payloads
func TestClientSignal(t *testing.T) {
	expectedSignalPayload := webwire.NewPayload(
		webwire.EncodingUtf8,
		[]byte("webwire_test_SIGNAL_payload"),
	)
	signalArrived := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	// Initialize webwire server given only the signal handler
	server := setupServer(
		t,
		&serverImpl{
			onSignal: func(
				_ context.Context,
				_ webwire.Connection,
				msg webwire.Message,
			) {

				// Verify signal payload
				comparePayload(
					t,
					"client signal",
					expectedSignalPayload,
					msg.Payload(),
				)

				// Synchronize, notify signal arrival
				signalArrived.Progress(1)
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
