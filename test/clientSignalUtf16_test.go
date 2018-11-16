package test

import (
	"context"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClientSignalUtf16 tests client-side signals with UTF16 encoded payloads
func TestClientSignalUtf16(t *testing.T) {
	testPayload := wwr.Payload{
		Encoding: wwr.EncodingUtf16,
		Data:     []byte{00, 115, 00, 97, 00, 109, 00, 112, 00, 108, 00, 101},
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
				assert.Equal(t, wwr.EncodingUtf16, msg.PayloadEncoding())
				assert.Equal(t, testPayload.Data, msg.Payload())

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
		wwr.Payload{
			Encoding: wwr.EncodingUtf16,
			Data: []byte{
				00, 115, 00, 97, 00, 109,
				00, 112, 00, 108, 00, 101,
			},
		},
	))

	// Synchronize, await signal arrival
	require.NoError(t, signalArrived.Wait(), "Signal wasn't processed")
}
