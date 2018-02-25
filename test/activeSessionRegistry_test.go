package test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestActiveSessionRegistry verifies that the session registry
// of currently active sessions is properly updated
func TestActiveSessionRegistry(t *testing.T) {
	authMsg := []byte("auth")
	logoutMsg := []byte("logout")

	// Initialize webwire server
	server := setupServer(
		t,
		webwire.Hooks{
			OnRequest: func(ctx context.Context) ([]byte, *webwire.Error) {
				// Extract request message and requesting client from the context
				msg := ctx.Value(webwire.MESSAGE).(webwire.Message)

				// Close session on logout
				if bytes.Equal(msg.Payload, logoutMsg) {
					if err := msg.Client.CloseSession(); err != nil {
						t.Errorf("Couldn't close session: %s", err)
					}
					return nil, nil
				}

				// Try to create a new session
				if err := msg.Client.CreateSession(nil); err != nil {
					return nil, &webwire.Error{
						Code:    "INTERNAL_ERROR",
						Message: fmt.Sprintf("Internal server error: %s", err),
					}
				}

				// Return the key of the newly created session
				return []byte(msg.Client.Session.Key), nil
			},
			// Define dummy hooks for sessions to be enabled on this server
			OnSessionCreated: func(_ *webwire.Client) error { return nil },
			OnSessionLookup:  func(_ string) (*webwire.Session, error) { return nil, nil },
			OnSessionClosed:  func(_ *webwire.Client) error { return nil },
		},
	)
	go server.Run()

	// Initialize client
	client := webwireClient.NewClient(
		server.Addr,
		webwireClient.Hooks{},
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)
	defer client.Close()

	// Send authentication request
	if _, err := client.Request(authMsg); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	activeSessionNumberBefore := server.SessionRegistry.Len()
	if activeSessionNumberBefore != 1 {
		t.Fatalf(
			"Unexpected active session number after authentication: %d",
			activeSessionNumberBefore,
		)
	}

	// Send logout request
	if _, err := client.Request(logoutMsg); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	activeSessionNumberAfter := server.SessionRegistry.Len()
	if activeSessionNumberAfter != 0 {
		t.Fatalf("Unexpected active session number after logout: %d", activeSessionNumberAfter)
	}
}
