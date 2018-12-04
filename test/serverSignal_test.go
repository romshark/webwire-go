package test

import (
	"sync"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServerSignal tests server-side signals
func TestServerSignal(t *testing.T) {
	expectedSignalPayload := wwr.Payload{
		Data: []byte("webwire_test_SERVER_SIGNAL_payload"),
	}
	signalProcessed := sync.WaitGroup{}
	signalProcessed.Add(1)

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			ClientConnected: func(
				_ wwr.ConnectionOptions,
				conn wwr.Connection,
			) {
				// Send signal
				assert.NoError(t, conn.Signal(
					nil,
					expectedSignalPayload,
				))
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		TestClientHooks{
			OnSignal: func(msg wwr.Message) {
				// Verify server signal payload
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

				// Synchronize, unlock main goroutine to pass the test case
				signalProcessed.Done()
			},
		},
	)

	// Synchronize, await signal arrival
	signalProcessed.Wait()
}
