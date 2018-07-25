package test

import (
	"context"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientSignalUtf16 tests client-side signals with UTF16 encoded payloads
func TestClientSignalUtf16(t *testing.T) {
	testPayload := webwire.NewPayload(
		webwire.EncodingUtf16,
		[]byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
	)
	verifyPayload := func(payload webwire.Payload) {
		payloadEncoding := payload.Encoding()
		if payloadEncoding != webwire.EncodingUtf16 {
			t.Errorf(
				"Unexpected payload encoding: %s",
				payloadEncoding.String(),
			)
		}
		testPayloadData := testPayload.Data()
		payloadData := payload.Data()
		if len(testPayloadData) != len(payloadData) {
			t.Errorf("Corrupt payload: %s", payloadData)
		}
		for i := 0; i < len(testPayloadData); i++ {
			if testPayloadData[i] != payloadData[i] {
				t.Errorf(
					"Corrupt payload, mismatching byte at position %d: %s",
					i,
					payloadData,
				)
				return
			}
		}
	}
	signalArrived := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	// Initialize webwire server given only the signal handler
	server := setupServer(
		t,
		&serverImpl{
			onSignal: func(
				_ context.Context,
				_ *webwire.Client,
				msg webwire.Message,
			) {
				verifyPayload(msg.Payload())

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
	err := client.connection.Signal("", webwire.NewPayload(
		webwire.EncodingUtf16,
		[]byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
	))
	if err != nil {
		t.Fatalf("Couldn't send signal: %s", err)
	}

	// Synchronize, await signal arrival
	if err := signalArrived.Wait(); err != nil {
		t.Fatal("Signal wasn't processed")
	}
}
