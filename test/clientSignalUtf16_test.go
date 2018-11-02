package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestClientSignalUtf16 tests client-side signals with UTF16 encoded payloads
func TestClientSignalUtf16(t *testing.T) {
	testPayload := wwr.NewPayload(
		wwr.EncodingUtf16,
		[]byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
	)
	verifyPayload := func(payload wwr.Payload) {
		assert.Equal(t, wwr.EncodingUtf16, payload.Encoding())
		assert.Equal(t, testPayload.Data(), payload.Data())
	}
	signalArrived := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	// Initialize webwire server given only the signal handler
	server := setupServer(
		t,
		&serverImpl{
			onSignal: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) {
				verifyPayload(msg.Payload())

				// Synchronize, notify signal arrival
				signalArrived.Progress(1)
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	require.NoError(t, client.connection.Connect())

	// Send signal
	require.NoError(t, client.connection.Signal(
		context.Background(),
		nil,
		wwr.NewPayload(
			wwr.EncodingUtf16,
			[]byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
		),
	))

	// Synchronize, await signal arrival
	require.NoError(t, signalArrived.Wait(), "Signal wasn't processed")
}
