package test

import (
	"context"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientOnSessionClosed tests the OnSessionClosed hook of the client
func TestClientOnSessionClosed(t *testing.T) {
	authenticated := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	hookCalled := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn webwire.Connection,
				_ webwire.Message,
			) (webwire.Payload, error) {
				// Try to create a new session
				if err := conn.CreateSession(nil); err != nil {
					return nil, err
				}

				go func() {
					// Wait until the authentication request is finished
					if err := authenticated.Wait(); err != nil {
						t.Errorf("Authentication timed out")
						return
					}

					// Close the session
					if err := conn.CloseSession(); err != nil {
						t.Errorf("Couldn't close session: %s", err)
					}
				}()

				return nil, nil
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
		callbackPoweredClientHooks{
			OnSessionClosed: func() {
				hookCalled.Progress(1)
			},
		},
	)

	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send authentication request and await reply
	if _, err := client.connection.Request(
		"login",
		webwire.NewPayload(webwire.EncodingBinary, []byte("credentials")),
	); err != nil {
		t.Fatalf("Request failed: %s", err)
	}
	authenticated.Progress(1)

	// Verify client session
	if err := hookCalled.Wait(); err != nil {
		t.Fatal("Hook not called")
	}
}
