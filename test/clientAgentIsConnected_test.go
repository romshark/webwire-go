package test

import (
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientAgentIsConnected tests the IsConnected method of the client agent
func TestClientAgentIsConnected(t *testing.T) {
	var clientAgent *wwr.Client
	clientAgentDefined := NewPending(1, 1*time.Second, true)
	clientDisconnected := NewPending(1, 1*time.Second, true)
	testerGoroutineFinished := NewPending(1, 1*time.Second, true)

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onClientConnected: func(newClt *wwr.Client) {
				if !newClt.IsConnected() {
					t.Errorf("Expected client agent to be connected")
				}
				clientAgent = newClt
				clientAgentDefined.Done()

				go func() {
					if err := clientDisconnected.Wait(); err != nil {
						t.Errorf("Client didn't disconnect")
					}

					if clientAgent.IsConnected() {
						t.Errorf("Expected client agent to be disconnected")
					}

					testerGoroutineFinished.Done()
				}()
			},
			onClientDisconnected: func(clt *wwr.Client) {
				if clientAgent.IsConnected() {
					t.Errorf("Expected client agent to be disconnected")
				}

				// Try to send a signal to a disconnected client and expect an error
				sigErr := clientAgent.Signal("", wwr.Payload{Data: []byte("testdata")})
				if _, isDisconnErr := sigErr.(wwr.DisconnectedErr); !isDisconnErr {
					t.Errorf(
						"Expected a DisconnectedErr, got: %s | %s",
						reflect.TypeOf(sigErr),
						sigErr,
					)
				}

				clientDisconnected.Done()
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := wwrclt.NewClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)

	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Wait for the client agent to be set by the OnClientConnected handler
	if err := clientAgentDefined.Wait(); err != nil {
		t.Fatalf("Tester goroutine didn't finish within 1 second")
	}

	if !clientAgent.IsConnected() {
		t.Fatalf("Expected client agent to be connected")
	}

	// Close the client connection and continue in the tester goroutine
	// spawned in the OnClientConnected handler of the server
	client.Close()

	// Wait for the tester goroutine to finish
	if err := testerGoroutineFinished.Wait(); err != nil {
		t.Fatalf("Tester goroutine didn't finish within 1 second")
	}
}
