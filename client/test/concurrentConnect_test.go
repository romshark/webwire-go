package client_test

import (
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	wwrtst "github.com/qbeon/webwire-go/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConcurrentConnect tests concurrent calling of Client.Connect
func TestConcurrentConnect(t *testing.T) {
	concurrentAccessors := 16
	finished := tmdwg.NewTimedWaitGroup(concurrentAccessors, 2*time.Second)

	// Initialize webwire server
	setup := wwrtst.SetupTestServer(
		t,
		&wwrtst.ServerImpl{},
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

	connect := func() {
		defer finished.Progress(1)
		assert.NoError(t, client.Connection.Connect())
	}

	for i := 0; i < concurrentAccessors; i++ {
		go connect()
	}

	require.NoError(t, finished.Wait(), "Expectation timed out")
}
