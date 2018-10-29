package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientConcurrentConnect verifies concurrent calling of client.Connect
// is properly synchronized and doesn't cause any data race
func TestClientConcurrentConnect(t *testing.T) {
	concurrentAccessors := 16
	finished := tmdwg.NewTimedWaitGroup(concurrentAccessors, 2*time.Second)

	// Initialize webwire server
	server := setupServer(t, &serverImpl{}, wwr.ServerOptions{})

	// Initialize client
	client := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)
	defer client.connection.Close()

	connect := func() {
		defer finished.Progress(1)
		assert.NoError(t, client.connection.Connect())
	}

	for i := 0; i < concurrentAccessors; i++ {
		go connect()
	}

	require.NoError(t, finished.Wait(), "Expectation timed out")
}
