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

// TestClientConcurrentRequest verifies concurrent calling of client.Request
// is properly synchronized and doesn't cause any data race
func TestClientConcurrentRequest(t *testing.T) {
	concurrentAccessors := 16
	finished := tmdwg.NewTimedWaitGroup(concurrentAccessors*2, 2*time.Second)

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				finished.Progress(1)
				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Address(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)
	defer client.connection.Close()

	require.NoError(t, client.connection.Connect())

	sendRequest := func() {
		defer finished.Progress(1)
		_, err := client.connection.Request(
			context.Background(),
			[]byte("sample"),
			wwr.Payload{
				Encoding: wwr.EncodingBinary,
				Data:     []byte("samplepayload"),
			},
		)
		assert.NoError(t, err)
	}

	for i := 0; i < concurrentAccessors; i++ {
		go sendRequest()
	}

	require.NoError(t, finished.Wait(), "Expectation timed out")
}
