package test

import (
	"context"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientOnSessionCreated tests the OnSessionCreated hook of the client
func TestClientOnSessionCreated(t *testing.T) {
	hookCalled := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	var createdSession *webwire.Session
	var sessionFromHook *webwire.Session

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				clt *webwire.Client,
				_ webwire.Message,
			) (webwire.Payload, error) {
				// Try to create a new session
				if err := clt.CreateSession(nil); err != nil {
					return nil, err
				}
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
			OnSessionCreated: func(newSession *webwire.Session) {
				sessionFromHook = newSession
				hookCalled.Progress(1)
			},
		},
	)
	defer client.connection.Close()

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

	createdSession = client.connection.Session()

	// Verify client session
	if err := hookCalled.Wait(); err != nil {
		t.Fatal("Hook not called")
	}

	// Compare the actual created session with the session received in the hook
	compareSessions(t, createdSession, sessionFromHook)
}
