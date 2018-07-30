package test

import (
	"reflect"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientConnIsConnected tests the IsActive method of a connection
func TestClientConnIsConnected(t *testing.T) {
	var clientConn wwr.Connection
	connectionReady := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	clientDisconnected := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	testerGoroutineFinished := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onClientConnected: func(newConn wwr.Connection) {
				if !newConn.IsActive() {
					t.Errorf("Expected connection to be active")
				}
				clientConn = newConn

				go func() {
					connectionReady.Progress(1)
					if err := clientDisconnected.Wait(); err != nil {
						t.Errorf("Client didn't disconnect")
					}

					if clientConn.IsActive() {
						t.Errorf("Expected connection to be inactive")
					}

					testerGoroutineFinished.Progress(1)
				}()
			},
			onClientDisconnected: func(_ wwr.Connection) {
				if clientConn.IsActive() {
					t.Errorf("Expected connection to be inactive")
				}

				// Try to send a signal to a inactive client and expect an error
				sigErr := clientConn.Signal("", wwr.NewPayload(
					wwr.EncodingBinary,
					[]byte("testdata"),
				))
				if _, isDisconnErr := sigErr.(wwr.DisconnectedErr); !isDisconnErr {
					t.Errorf(
						"Expected a DisconnectedErr, got: %s | %s",
						reflect.TypeOf(sigErr),
						sigErr,
					)
				}

				clientDisconnected.Progress(1)
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
			Autoconnect:           wwr.Disabled,
		},
		callbackPoweredClientHooks{},
	)

	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Wait for the connection to be set by the OnClientConnected handler
	if err := connectionReady.Wait(); err != nil {
		t.Fatalf("Connection not ready after 1 second")
	}

	if !clientConn.IsActive() {
		t.Fatalf("Expected connection to be active")
	}

	// Close the client connection and continue in the tester goroutine
	// spawned in the OnClientConnected handler of the server
	client.connection.Close()

	// Wait for the tester goroutine to finish
	if err := testerGoroutineFinished.Wait(); err != nil {
		t.Fatalf("Tester goroutine didn't finish within 1 second")
	}
}
