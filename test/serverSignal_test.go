package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestServerSignal tests server-side signals
func TestServerSignal(t *testing.T) {
	expectedSignalPayload := wwr.NewPayload(
		wwr.EncodingBinary,
		[]byte("webwire_test_SERVER_SIGNAL_payload"),
	)
	signalProcessed := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onClientConnected: func(conn wwr.Connection) {
				// Send signal
				assert.NoError(t, conn.Signal(
					"",
					expectedSignalPayload,
				))
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
		callbackPoweredClientHooks{
			OnSignal: func(signalMessage wwr.Message) {
				// Verify server signal payload
				comparePayload(
					t,
					expectedSignalPayload,
					signalMessage.Payload(),
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
