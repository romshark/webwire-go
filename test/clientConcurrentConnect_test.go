package test

import (
	"os"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientConcurrentConnect verifies concurrent calling of client.Connect
// is properly synchronized and doesn't cause any data race
func TestClientConcurrentConnect(t *testing.T) {
	var concurrentAccessors uint32 = 16
	finished := NewPending(concurrentAccessors, 2*time.Second, true)

	// Initialize webwire server
	server := setupServer(t, webwire.Hooks{})
	go server.Run()

	// Initialize client
	client := webwireClient.NewClient(
		server.Addr,
		webwireClient.Hooks{},
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)
	defer client.Close()

	connect := func() {
		defer finished.Done()
		if err := client.Connect(); err != nil {
			t.Errorf("Connect failed: %s", err)
		}
	}

	for i := uint32(0); i < concurrentAccessors; i++ {
		go connect()
	}

	if err := finished.Wait(); err != nil {
		t.Fatal("Expectation timed out")
	}
}
