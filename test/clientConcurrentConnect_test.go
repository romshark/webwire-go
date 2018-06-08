package test

import (
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientConcurrentConnect verifies concurrent calling of client.Connect
// is properly synchronized and doesn't cause any data race
func TestClientConcurrentConnect(t *testing.T) {
	concurrentAccessors := 16
	finished := tmdwg.NewTimedWaitGroup(concurrentAccessors, 2*time.Second)

	// Initialize webwire server
	server := setupServer(t, &serverImpl{}, webwire.ServerOptions{})

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)
	defer client.connection.Close()

	connect := func() {
		defer finished.Progress(1)
		if err := client.connection.Connect(); err != nil {
			t.Errorf("Connect failed: %s", err)
		}
	}

	for i := 0; i < concurrentAccessors; i++ {
		go connect()
	}

	if err := finished.Wait(); err != nil {
		t.Fatal("Expectation timed out")
	}
}
