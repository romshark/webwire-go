package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientOnSessionCreated verifies the OnSessionCreated hook of the client is called properly.
func TestClientOnSessionCreated(t *testing.T) {
	hookCalled := NewPending(1, 1*time.Second, true)
	var createdSession *webwire.Session
	var sessionFromHook *webwire.Session

	// Initialize webwire server
	_, addr := setupServer(
		t,
		webwire.ServerOptions{
			Hooks: webwire.Hooks{
				OnRequest: func(ctx context.Context) (webwire.Payload, error) {
					// Extract request message and requesting client from the context
					msg := ctx.Value(webwire.Msg).(webwire.Message)

					// Try to create a new session
					if err := msg.Client.CreateSession(nil); err != nil {
						return webwire.Payload{}, err
					}
					return webwire.Payload{}, nil
				},
				// Define dummy hooks to enable sessions on this server
				OnSessionCreated: func(_ *webwire.Client) error { return nil },
				OnSessionLookup:  func(_ string) (*webwire.Session, error) { return nil, nil },
				OnSessionClosed:  func(_ *webwire.Client) error { return nil },
			},
		},
	)

	// Initialize client
	client := webwireClient.NewClient(
		addr,
		webwireClient.Options{
			Hooks: webwireClient.Hooks{
				OnSessionCreated: func(newSession *webwire.Session) {
					sessionFromHook = newSession
					hookCalled.Done()
				},
			},
			DefaultRequestTimeout: 2 * time.Second,
		},
	)
	defer client.Close()

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

	tmp := client.Session()
	createdSession = &tmp

	// Verify client session
	if err := hookCalled.Wait(); err != nil {
		t.Fatal("Hook not called")
	}

	// Compare the actual created session with the session received in the hook
	compareSessions(t, createdSession, sessionFromHook)
}
