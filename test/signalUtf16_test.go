package test

import (
	"context"
	"sync"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/assert"
)

// TestSignalUtf16 tests client-side signals with UTF16 encoded payloads
func TestSignalUtf16(t *testing.T) {
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
				assert.Equal(t, wwr.EncodingUtf16, msg.PayloadEncoding())
				assert.Equal(t, []byte{
					00, 115, 00, 97, 00, 109,
					00, 112, 00, 108, 00, 101,
				}, msg.Payload())

				// Synchronize, notify signal arrival
				signalArrived.Done()
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	signal(t, sock, []byte("utf16_sig"), payload.Payload{
		Encoding: wwr.EncodingUtf16,
		Data: []byte{
			00, 115, 00, 97, 00, 109,
			00, 112, 00, 108, 00, 101,
		},
	})

	// Synchronize, await signal arrival
	signalArrived.Wait()
}
