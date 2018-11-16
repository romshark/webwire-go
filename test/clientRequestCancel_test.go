package test

import (
	"context"
	"testing"
	"time"

	"github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClientRequestCancel tests canceling of fired requests
func TestClientRequestCancel(t *testing.T) {
	requestFinished := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	// Initialize webwire server given only the request
	setup := setupTestServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				time.Sleep(2 * time.Second)
				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.newClient(
		wwrclt.Options{
			DefaultRequestTimeout: 5 * time.Second,
		},
		nil, // Use the default transport implementation
		testClientHooks{},
	)

	require.NoError(t, client.connection.Connect())

	cancelableCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Send request and await reply
	go func() {
		reply, err := client.connection.Request(
			cancelableCtx,
			[]byte("test"),
			wwr.Payload{},
		)
		assert.Error(t, err, "Expected a canceled-error")
		assert.Nil(t, reply)
		assert.IsType(t, wwr.CanceledErr{}, err)
		assert.True(t, wwr.IsCanceledErr(err))
		assert.False(t, wwr.IsTimeoutErr(err))
		requestFinished.Progress(1)
	}()

	// Cancel the context some time after sending the request
	time.Sleep(10 * time.Millisecond)
	cancel()

	// Wait for the requestor goroutine to finish
	require.NoError(t, requestFinished.Wait(), "Test timed out")
}
