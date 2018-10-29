package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientConcurrentSignal verifies concurrent calling of client.Signal
// is properly synchronized and doesn't cause any data race
func TestClientConcurrentSignal(t *testing.T) {
	concurrentAccessors := 16
	finished := tmdwg.NewTimedWaitGroup(concurrentAccessors*2, 2*time.Second)

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onSignal: func(
				_ context.Context,
				_ wwr.Connection,
				_ wwr.Message,
			) {
				finished.Progress(1)
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
		callbackPoweredClientHooks{},
	)
	defer client.connection.Close()

	require.NoError(t, client.connection.Connect())

	sendSignal := func() {
		defer finished.Progress(1)
		assert.NoError(t, client.connection.Signal(
			"sample",
			wwr.NewPayload(wwr.EncodingBinary, []byte("samplepayload")),
		))
	}

	for i := 0; i < concurrentAccessors; i++ {
		go sendSignal()
	}

	require.NoError(t, finished.Wait(), "Expectation timed out")
}
