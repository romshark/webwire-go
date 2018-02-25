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

// TestClientOnSessionCreated verifies the OnSessionCreated hook of the client is called properly.
func TestClientOnSessionCreated(t *testing.T) {
	var hook sync.WaitGroup
	hook.Add(1)
	var createdSession *webwire.Session
	var sessionFromHook *webwire.Session

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
			OnSessionCreated: func(newSession *webwire.Session) {
				sessionFromHook = newSession
				hook.Done()
			},
		},
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)
	defer client.Close()

	// Send request and await reply
	_, err := client.Request([]byte("credentials"))
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	tmp := client.Session()
	createdSession = &tmp

	// Verify client session
	hook.Wait()

	// Compare the actual created session with the session received in the hook
	compareSessions(t, createdSession, sessionFromHook)
}
