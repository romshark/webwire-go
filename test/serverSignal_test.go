package test

import (
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
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
	signalProcessed := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onClientConnected: func(conn wwr.Connection) {
				// Send signal
				assert.NoError(t, conn.Signal(
					nil,
					expectedSignalPayload,
				))
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
		callbackPoweredClientHooks{
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
				signalProcessed.Progress(1)
			},
		},
	)
	defer client.connection.Close()

	// Connect client
	require.NoError(t, client.connection.Connect())

	// Synchronize, await signal arrival
	require.NoError(t,
		signalProcessed.Wait(),
		"Server signal didn't arrive",
	)
}
