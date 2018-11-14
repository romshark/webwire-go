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

// TestClientConcurrentSignal verifies concurrent calling of client.Signal
// is properly synchronized and doesn't cause any data race
func TestClientConcurrentSignal(t *testing.T) {
	concurrentAccessors := 16
	finished := tmdwg.NewTimedWaitGroup(concurrentAccessors*2, 2*time.Second)

	// Initialize webwire server
	setup := setupTestServer(
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
	defer client.connection.Close()

	require.NoError(t, client.connection.Connect())

	sendSignal := func() {
		defer finished.Progress(1)
		assert.NoError(t, client.connection.Signal(
			context.Background(),
			[]byte("sample"),
			wwr.Payload{
				Encoding: wwr.EncodingBinary,
				Data:     []byte("samplepayload"),
			},
		))
	}

	for i := 0; i < concurrentAccessors; i++ {
		go sendSignal()
	}

	require.NoError(t, finished.Wait(), "Expectation timed out")
}
