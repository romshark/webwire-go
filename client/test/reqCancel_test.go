package client_test

import (
	"context"
	"sync"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	wwrtst "github.com/qbeon/webwire-go/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReqCancel tests canceling of fired requests
func TestReqCancel(t *testing.T) {
	requestFinished := sync.WaitGroup{}
	requestFinished.Add(1)

	// Initialize webwire server given only the request
	setup := wwrtst.SetupTestServer(
		t,
		&wwrtst.ServerImpl{
			Request: func(
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
	client := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 5 * time.Second,
		},
		nil, // Use the default transport implementation
		wwrtst.TestClientHooks{},
	)

	require.NoError(t, client.Connection.Connect())

	cancelableCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Send request and await reply
	go func() {
		reply, err := client.Connection.Request(
			cancelableCtx,
			[]byte("test"),
			wwr.Payload{},
		)
		assert.Error(t, err, "Expected a canceled-error")
		assert.Nil(t, reply)
		assert.IsType(t, wwr.CanceledErr{}, err)
		assert.True(t, wwr.IsCanceledErr(err))
		assert.False(t, wwr.IsTimeoutErr(err))
		requestFinished.Done()
	}()

	// Cancel the context some time after sending the request
	time.Sleep(10 * time.Millisecond)
	cancel()

	// Wait for the requestor goroutine to finish
	requestFinished.Wait()
}
