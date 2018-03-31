package test

import (
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestServerSignal tests server-side signals
func TestServerSignal(t *testing.T) {
	expectedSignalPayload := webwire.Payload{
		Data: []byte("webwire_test_SERVER_SIGNAL_payload"),
	}
	var addr string
	serverSignalArrived := NewPending(1, 1*time.Second, true)
	initClient := make(chan bool, 1)
	sendSignal := make(chan bool, 1)

	// Initialize webwire server
	go func() {
		server := setupServer(
			t,
			webwire.ServerOptions{
				Hooks: webwire.Hooks{
					OnClientConnected: func(client *webwire.Client) {
						// Send signal
						if err := client.Signal(
							"",
							expectedSignalPayload,
						); err != nil {
							t.Fatalf("Couldn't send signal to client: %s", err)
						}
					},
				},
			},
		)
		addr = server.Addr().String()

		// Synchronize, initialize client
		initClient <- true

		// Synchronize, wait for the client to launch
		// and require the signal to be sent
		<-sendSignal
	}()

	// Synchronize, await server initialization
	<-initClient

	// Initialize client
	client := webwireClient.NewClient(
		addr,
		webwireClient.Options{
			Hooks: webwireClient.Hooks{
				OnServerSignal: func(signalPayload webwire.Payload) {
					// Verify server signal payload
					comparePayload(
						t,
						"server signal",
						expectedSignalPayload,
						signalPayload,
					)

					// Synchronize, unlock main goroutine to pass the test case
					serverSignalArrived.Done()
				},
			},
			DefaultRequestTimeout: 2 * time.Second,
		},
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
	if err := serverSignalArrived.Wait(); err != nil {
		t.Fatal("Server signal didn't arrive")
	}
}
