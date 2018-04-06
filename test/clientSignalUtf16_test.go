package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientSignalUtf16 tests client-side signals with UTF16 encoded payloads
func TestClientSignalUtf16(t *testing.T) {
	testPayload := webwire.Payload{
		Encoding: webwire.EncodingUtf16,
		Data:     []byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
	}
	verifyPayload := func(payload webwire.Payload) {
		if payload.Encoding != webwire.EncodingUtf16 {
			t.Errorf(
				"Unexpected payload encoding: %s",
				payload.Encoding.String(),
			)
		}
		if len(testPayload.Data) != len(payload.Data) {
			t.Errorf("Corrupt payload: %s", payload.Data)
		}
		for i := 0; i < len(testPayload.Data); i++ {
			if testPayload.Data[i] != payload.Data[i] {
				t.Errorf(
					"Corrupt payload, mismatching byte at position %d: %s",
					i,
					payload.Data,
				)
				return
			}
		}
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
				verifyPayload(msg.Payload)

				// Synchronize, notify signal arrival
				signalArrived.Done()
			},
		},
		webwire.ServerOptions{},
	)

	// Initialize client
	client := webwireClient.NewClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)

	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send signal
	err := client.Signal("", webwire.Payload{
		Encoding: webwire.EncodingUtf16,
		Data:     []byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
	})
	if err != nil {
		t.Fatalf("Couldn't send signal: %s", err)
	}

	// Synchronize, await signal arrival
	if err := signalArrived.Wait(); err != nil {
		t.Fatal("Signal wasn't processed")
	}
}
