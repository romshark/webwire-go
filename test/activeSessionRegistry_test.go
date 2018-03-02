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

// TestActiveSessionRegistry verifies that the session registry
// of currently active sessions is properly updated
func TestActiveSessionRegistry(t *testing.T) {
	// Initialize webwire server
	srv, addr := setupServer(
		t,
		webwire.Hooks{
			OnRequest: func(ctx context.Context) (webwire.Payload, *webwire.Error) {
				// Extract request message and requesting client from the context
				msg := ctx.Value(webwire.Msg).(webwire.Message)

				// Close session on logout
				if msg.Name == "logout" {
					if err := msg.Client.CloseSession(); err != nil {
						t.Errorf("Couldn't close session: %s", err)
					}
					return webwire.Payload{}, nil
				}

				// Try to create a new session
				if err := msg.Client.CreateSession(nil); err != nil {
					return webwire.Payload{}, &webwire.Error{
						Code:    "INTERNAL_ERROR",
						Message: fmt.Sprintf("Internal server error: %s", err),
					}
				}

				// Return the key of the newly created session (use default binary encoding)
				return webwire.Payload{
					Data: []byte(msg.Client.Session.Key),
				}, nil
			},
			// Define dummy hooks for sessions to be enabled on this server
			OnSessionCreated: func(_ *webwire.Client) error { return nil },
			OnSessionLookup:  func(_ string) (*webwire.Session, error) { return nil, nil },
			OnSessionClosed:  func(_ *webwire.Client) error { return nil },
		},
	)

	// Initialize client
	client := webwireClient.NewClient(
		addr,
		webwireClient.Hooks{},
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)
	defer client.Close()

	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send authentication request
	if _, err := client.Request(
		"login",
		webwire.Payload{
			Encoding: webwire.EncodingUtf8,
			Data:     []byte("nothing"),
		},
	); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	activeSessionNumberBefore := srv.SessionRegistry.Len()
	if activeSessionNumberBefore != 1 {
		t.Fatalf(
			"Unexpected active session number after authentication: %d",
			activeSessionNumberBefore,
		)
	}

	// Send logout request
	if _, err := client.Request(
		"logout",
		webwire.Payload{
			Encoding: webwire.EncodingUtf8,
			Data:     []byte("nothing"),
		},
	); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	activeSessionNumberAfter := srv.SessionRegistry.Len()
	if activeSessionNumberAfter != 0 {
		t.Fatalf("Unexpected active session number after logout: %d", activeSessionNumberAfter)
	}
}
