package client_test

import (
	"context"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	wwrtst "github.com/qbeon/webwire-go/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConcurrentRequest tests concurrent calling of Client.Request
func TestConcurrentRequest(t *testing.T) {
	concurrentAccessors := 16
	finished := tmdwg.NewTimedWaitGroup(concurrentAccessors*2, 2*time.Second)

	// Initialize webwire server
	setup := wwrtst.SetupTestServer(
		t,
		&wwrtst.ServerImpl{
			Request: func(
				_ context.Context,
				_ wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				finished.Progress(1)
				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		wwrtst.TestClientHooks{},
	)
	defer client.Connection.Close()

	require.NoError(t, client.Connection.Connect())

	sendRequest := func() {
		defer finished.Progress(1)
		_, err := client.Connection.Request(
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
