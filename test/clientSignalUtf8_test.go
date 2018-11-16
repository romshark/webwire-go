package test

import (
	"context"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestClientSignal tests client-side signals with UTF8 encoded payloads
func TestClientSignal(t *testing.T) {
	expectedSignalPayload := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     []byte("webwire_test_SIGNAL_payload"),
	}
	signalArrived := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	// Initialize webwire server given only the signal handler
	setup := setupTestServer(
		t,
		&serverImpl{
			onSignal: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) {
				// Verify signal payload
				require.Equal(
					t,
					expectedSignalPayload.Encoding,
					msg.PayloadEncoding(),
				)
				require.Equal(
					t,
					expectedSignalPayload.Data,
					msg.Payload(),
				)

				// Synchronize, notify signal arrival
				signalArrived.Progress(1)
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.newClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		testClientHooks{},
	)

	require.NoError(t, client.connection.Connect())

	// Send signal
	require.NoError(t, client.connection.Signal(
		context.Background(),
		nil,
		expectedSignalPayload,
	))

	// Synchronize, await signal arrival
	require.NoError(t, signalArrived.Wait(), "Signal wasn't processed")
}
