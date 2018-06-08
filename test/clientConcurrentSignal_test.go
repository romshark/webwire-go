package test

import (
	"context"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientConcurrentSignal verifies concurrent calling of client.Signal
// is properly synchronized and doesn't cause any data race
func TestClientConcurrentSignal(t *testing.T) {
	concurrentAccessors := 16
	finished := tmdwg.NewTimedWaitGroup(concurrentAccessors*2, 2*time.Second)

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onSignal: func(
				_ context.Context,
				_ *webwire.Client,
				_ *webwire.Message,
			) {
				finished.Progress(1)
			},
		},
		webwire.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)
	defer client.connection.Close()

	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	sendSignal := func() {
		defer finished.Progress(1)
		if err := client.connection.Signal(
			"sample",
			webwire.Payload{Data: []byte("samplepayload")},
		); err != nil {
			t.Errorf("Request failed: %s", err)
		}
	}

	for i := 0; i < concurrentAccessors; i++ {
		go sendSignal()
	}

	if err := finished.Wait(); err != nil {
		t.Fatal("Expectation timed out")
	}
}
