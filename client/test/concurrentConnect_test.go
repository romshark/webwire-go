package client_test

import (
	"sync"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	wwrtst "github.com/qbeon/webwire-go/test"
	"github.com/stretchr/testify/assert"
)

// TestConcurrentConnect tests concurrent calling of Client.Connect
func TestConcurrentConnect(t *testing.T) {
	concurrentAccessors := 16
	finished := sync.WaitGroup{}
	finished.Add(concurrentAccessors)

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
		defer finished.Done()
		assert.NoError(t, client.Connection.Connect())
	}

	for i := 0; i < concurrentAccessors; i++ {
		go connect()
	}

	finished.Wait()
}
