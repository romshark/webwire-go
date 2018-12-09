package test

import (
	"context"
	"sync"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/require"
)

// TestSignalUtf8 tests client-side signals with UTF8 encoded payloads
func TestSignalUtf8(t *testing.T) {
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
				require.Equal(t, wwr.EncodingUtf8, msg.PayloadEncoding())
				require.Equal(t, []byte("üникод"), msg.Payload())

				// Synchronize, notify signal arrival
				signalArrived.Done()
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	signal(t, sock, []byte("sig_utf8"), payload.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     []byte("üникод"),
	})

	// Synchronize, await signal arrival
	signalArrived.Wait()
}
