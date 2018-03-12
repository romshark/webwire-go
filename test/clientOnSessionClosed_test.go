package test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientOnSessionClosed verifies the OnSessionClosed hook of the client is called properly.
func TestClientOnSessionClosed(t *testing.T) {
	authenticated := NewPending(1, 1*time.Second, true)
	hookCalled := NewPending(1, 1*time.Second, true)

	// Initialize webwire server
	_, addr := setupServer(
		t,
		webwire.Hooks{
			OnRequest: func(ctx context.Context) (webwire.Payload, error) {
				// Extract request message and requesting client from the context
				msg := ctx.Value(webwire.Msg).(webwire.Message)

				// Try to create a new session
				if err := msg.Client.CreateSession(nil); err != nil {
					return webwire.Payload{}, webwire.Error{
						Code:    "INTERNAL_ERROR",
						Message: fmt.Sprintf("Internal server error: %s", err),
					}
				}

				go func() {
					// Wait until the authentication request is finished
					if err := authenticated.Wait(); err != nil {
						t.Errorf("Authentication timed out")
						return
					}

					// Close the session
					if err := msg.Client.CloseSession(); err != nil {
						t.Errorf("Couldn't close session: %s", err)
					}
				}()

				return webwire.Payload{}, nil
			},
			// Define dummy hooks to enable sessions on this server
			OnSessionCreated: func(_ *webwire.Client) error { return nil },
			OnSessionLookup:  func(_ string) (*webwire.Session, error) { return nil, nil },
			OnSessionClosed:  func(_ *webwire.Client) error { return nil },
		},
	)

	// Initialize client
	client := webwireClient.NewClient(
		addr,
		webwireClient.Hooks{
			OnSessionClosed: func() {
				hookCalled.Done()
			},
		},
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)

	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send authentication request and await reply
	if _, err := client.Request(
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
