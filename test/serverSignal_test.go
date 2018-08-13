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
	var addr string
	serverSignalArrived := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	initClient := make(chan bool, 1)
	sendSignal := make(chan bool, 1)

	// Initialize webwire server
	go func() {
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
		addr = server.Addr().String()

		// Synchronize, initialize client
		initClient <- true

		// Synchronize, wait for the client to launch
		// and require the signal to be sent
		<-sendSignal
	}()

	// Synchronize, await server initialization
	<-initClient

	// Initialize client
	client := newCallbackPoweredClient(
		addr,
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{
			OnSignal: func(signalPayload wwr.Payload) {
				// Verify server signal payload
				comparePayload(t, expectedSignalPayload, signalPayload)

				// Synchronize, unlock main goroutine to pass the test case
				serverSignalArrived.Progress(1)
			},
		},
	)
	defer client.connection.Close()

	// Connect client
	require.NoError(t, client.connection.Connect())

	// Synchronize, notify the server the client was initialized
	// and request the signal
	sendSignal <- true

	// Synchronize, await signal arrival
	require.NoError(t,
		serverSignalArrived.Wait(),
		"Server signal didn't arrive",
	)
}
