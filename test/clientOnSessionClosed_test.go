package test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientOnSessionClosed verifies the OnSessionClosed hook of the client is called properly.
func TestClientOnSessionClosed(t *testing.T) {
	var authentication sync.WaitGroup
	var hook sync.WaitGroup
	authentication.Add(1)
	hook.Add(1)

	// Initialize webwire server
	server := setupServer(
		t,
		webwire.Hooks{
			OnRequest: func(ctx context.Context) ([]byte, *webwire.Error) {
				// Extract request message and requesting client from the context
				msg := ctx.Value(webwire.MESSAGE).(webwire.Message)

				// Try to create a new session
				if err := msg.Client.CreateSession(nil); err != nil {
					return nil, &webwire.Error{
						Code:    "INTERNAL_ERROR",
						Message: fmt.Sprintf("Internal server error: %s", err),
					}
				}

				go func() {
					// Wait until the authentication request is finished
					authentication.Wait()

					// Close the session
					if err := msg.Client.CloseSession(); err != nil {
						t.Errorf("Couldn't close session: %s", err)
					}
				}()

				return nil, nil
			},
			// Define dummy hooks to enable sessions on this server
			OnSessionCreated: func(_ *webwire.Client) error { return nil },
			OnSessionLookup:  func(_ string) (*webwire.Session, error) { return nil, nil },
			OnSessionClosed:  func(_ *webwire.Client) error { return nil },
		},
	)
	go server.Run()

	// Initialize client
	client := webwireClient.NewClient(
		server.Addr,
		webwireClient.Hooks{
			OnSessionClosed: func() {
				hook.Done()
			},
		},
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)

	// Send request and await reply
	_, err := client.Request([]byte("credentials"))
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}
	authentication.Done()

	// Verify client session
	hook.Wait()
}
