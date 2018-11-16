package test

import (
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClientConcurrentConnect verifies concurrent calling of client.Connect
// is properly synchronized and doesn't cause any data race
func TestClientConcurrentConnect(t *testing.T) {
	concurrentAccessors := 16
	finished := tmdwg.NewTimedWaitGroup(concurrentAccessors, 2*time.Second)

	// Initialize webwire server
	setup := setupTestServer(
		t,
		&serverImpl{},
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

	connect := func() {
		defer finished.Progress(1)
		assert.NoError(t, client.connection.Connect())
	}

	for i := 0; i < concurrentAccessors; i++ {
		go connect()
	}

	require.NoError(t, finished.Wait(), "Expectation timed out")
}
