package test

import (
	"os"
	"sync"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwire_client "github.com/qbeon/webwire-go/client"
)

// TestServerSignal verifies the server is connectable
// and sends signals correctly
func TestServerSignal(t *testing.T) {
	expectedSignalPayload := []byte("webwire_test_SERVER_SIGNAL_payload")
	var addr string
	var server *webwire.Server
	var finish sync.WaitGroup
	finish.Add(1)
	initClient := make(chan bool, 1)
	sendSignal := make(chan bool, 1)

	// Initialize webwire server
	go func() {
		server = setupServer(t, func(client *webwire.Client) {

			// Verify client is listed
			/*
				if server.ClientsNum() != 1 {
					finish.Done()
					t.Fatalf(
						"Unexpected list of connected clients (%d), "+
							"expected 1 client to be connected",
						server.ClientsNum(),
					)
				}
			*/

			// Send signal
			if err := client.Signal(expectedSignalPayload); err != nil {
				t.Fatalf("Couldn't send signal to client: %s", err)
			}
		}, nil, nil, nil, nil, nil, nil, nil)
		go server.Run()
		addr = server.Addr

		// Synchronize, initialize client
		initClient <- true

		// Synchronize, wait for the client to launch
		// and require the signal to be sent
		<-sendSignal
	}()

	// Synchronize, await server initialization
	<-initClient

	// Initialize client
	client := webwire_client.NewClient(
		addr,
		func(signalPayload []byte) {
			// Verify server signal payload
			comparePayload(
				t,
				"server signal",
				expectedSignalPayload,
				signalPayload,
			)

			// Synchronize, unlock main goroutine to pass the test case
			finish.Done()
		},
		nil, nil,
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)
	defer client.Close()

	// Connect client
	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect client: %s", err)
	}

	// Synchronize, notify the server the client was initialized
	// and request the signal
	sendSignal <- true

	// Synchronize, await signal arrival
	finish.Wait()
}
