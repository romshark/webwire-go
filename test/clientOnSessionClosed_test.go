package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientOnSessionClosed tests the OnSessionClosed hook of the client
func TestClientOnSessionClosed(t *testing.T) {
	authenticated := newPending(1, 1*time.Second, true)
	hookCalled := newPending(1, 1*time.Second, true)

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				clt *webwire.Client,
				_ *webwire.Message,
			) (webwire.Payload, error) {
				// Try to create a new session
				if err := clt.CreateSession(nil); err != nil {
					return webwire.Payload{}, err
				}

				go func() {
					// Wait until the authentication request is finished
					if err := authenticated.Wait(); err != nil {
						t.Errorf("Authentication timed out")
						return
					}

					// Close the session
					if err := clt.CloseSession(); err != nil {
						t.Errorf("Couldn't close session: %s", err)
					}
				}()

				return webwire.Payload{}, nil
			},
		},
		webwire.ServerOptions{
			SessionsEnabled: true,
		},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil,
		func() {
			hookCalled.Done()
		},
		nil, nil,
	)

	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send authentication request and await reply
	if _, err := client.connection.Request(
		"login",
		webwire.Payload{Data: []byte("credentials")},
	); err != nil {
		t.Fatalf("Request failed: %s", err)
	}
	authenticated.Done()

	// Verify client session
	if err := hookCalled.Wait(); err != nil {
		t.Fatal("Hook not called")
	}
}
