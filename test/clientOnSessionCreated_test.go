package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientOnSessionCreated tests the OnSessionCreated hook of the client
func TestClientOnSessionCreated(t *testing.T) {
	hookCalled := NewPending(1, 1*time.Second, true)
	var createdSession *webwire.Session
	var sessionFromHook *webwire.Session

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(ctx context.Context) (webwire.Payload, error) {
				// Extract request message
				// and requesting client from the context
				msg := ctx.Value(webwire.Msg).(webwire.Message)

				// Try to create a new session
				if err := msg.Client.CreateSession(nil); err != nil {
					return webwire.Payload{}, err
				}
				return webwire.Payload{}, nil
			},
		},
		webwire.ServerOptions{
			SessionsEnabled: true,
		},
	)

	// Initialize client
	client := webwireClient.NewClient(
		server.Addr().String(),
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
