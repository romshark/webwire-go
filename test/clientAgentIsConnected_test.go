package test

import (
	"reflect"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientAgentIsConnected tests the IsConnected method of the client agent
func TestClientAgentIsConnected(t *testing.T) {
	var clientAgent *wwr.Client
	clientAgentReady := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	clientDisconnected := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	testerGoroutineFinished := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onClientConnected: func(newClt *wwr.Client) {
				if !newClt.IsConnected() {
					t.Errorf("Expected client agent to be connected")
				}
				clientAgent = newClt

				go func() {
					clientAgentReady.Progress(1)
					if err := clientDisconnected.Wait(); err != nil {
						t.Errorf("Client didn't disconnect")
					}

					if clientAgent.IsConnected() {
						t.Errorf("Expected client agent to be disconnected")
					}

					testerGoroutineFinished.Progress(1)
				}()
			},
			onClientDisconnected: func(clt *wwr.Client) {
				if clientAgent.IsConnected() {
					t.Errorf("Expected client agent to be disconnected")
				}

				// Try to send a signal to a disconnected client and expect an error
				sigErr := clientAgent.Signal("", wwr.NewPayload(
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

	// Wait for the client agent to be set by the OnClientConnected handler
	if err := clientAgentReady.Wait(); err != nil {
		t.Fatalf("Client agent not ready after 1 second")
	}

	if !clientAgent.IsConnected() {
		t.Fatalf("Expected client agent to be connected")
	}

	// Close the client connection and continue in the tester goroutine
	// spawned in the OnClientConnected handler of the server
	client.connection.Close()

	// Wait for the tester goroutine to finish
	if err := testerGoroutineFinished.Wait(); err != nil {
		t.Fatalf("Tester goroutine didn't finish within 1 second")
	}
}
