package test

import (
	"context"
	"sync"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestSignalUtf8 tests client-side signals with UTF8 encoded payloads
func TestSignalUtf8(t *testing.T) {
	expectedSignalPayload := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     []byte("webwire_test_SIGNAL_payload"),
	}
	signalArrived := sync.WaitGroup{}
	signalArrived.Add(1)

	// Initialize webwire server given only the signal handler
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Signal: func(
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
				signalArrived.Done()
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		TestClientHooks{},
	)

	require.NoError(t, client.Connection.Connect())

	// Send signal
	require.NoError(t, client.Connection.Signal(
		context.Background(),
		nil,
		expectedSignalPayload,
	))

	// Synchronize, await signal arrival
	signalArrived.Wait()
}
