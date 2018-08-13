package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientSignal tests client-side signals with UTF8 encoded payloads
func TestClientSignal(t *testing.T) {
	expectedSignalPayload := wwr.NewPayload(
		wwr.EncodingUtf8,
		[]byte("webwire_test_SIGNAL_payload"),
	)
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

				// Verify signal payload
				comparePayload(t, expectedSignalPayload, msg.Payload())

				// Synchronize, notify signal arrival
				signalArrived.Progress(1)
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	require.NoError(t, client.connection.Connect())

	// Send signal
	require.NoError(t, client.connection.Signal("", expectedSignalPayload))

	// Synchronize, await signal arrival
	require.NoError(t, signalArrived.Wait(), "Signal wasn't processed")
}
