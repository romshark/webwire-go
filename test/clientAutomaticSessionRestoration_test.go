package test

import (
	"context"
	"fmt"
	"testing"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientAutomaticSessionRestoration verifies automatic session restoration
// on connection establishment
func TestClientAutomaticSessionRestoration(t *testing.T) {
	sessionStorage := make(map[string]*webwire.Session)

	currentStep := 1
	var createdSession *webwire.Session

	// Initialize webwire server
	_, addr := setupServer(
		t,
		webwire.ServerOptions{
			SessionsEnabled: true,
			Hooks: webwire.Hooks{
				OnRequest: func(ctx context.Context) (webwire.Payload, error) {
					// Extract request message and requesting client from the context
					msg := ctx.Value(webwire.Msg).(webwire.Message)

					if currentStep == 2 {
						// Expect the session to have been automatically restored
						compareSessions(t, createdSession, msg.Client.Session)
						return webwire.Payload{}, nil
					}

					// Try to create a new session
					if err := msg.Client.CreateSession(nil); err != nil {
						return webwire.Payload{}, err
					}

					// Return the key of the newly created session
					return webwire.Payload{}, nil
				},
				// Permanently store the session
				OnSessionCreated: func(client *webwire.Client) error {
					sessionStorage[client.Session.Key] = client.Session
					return nil
				},
				// Find session by key
				OnSessionLookup: func(key string) (*webwire.Session, error) {
					// Expect the key of the created session to be looked up
					if key != createdSession.Key {
						err := fmt.Errorf(
							"Expected and looked up session keys differ: %s | %s",
							createdSession.Key,
							key,
						)
						t.Fatalf("Session lookup mismatch: %s", err)
						return nil, err
					}

					if session, exists := sessionStorage[key]; exists {
						return session, nil
					}

					// Expect the session to be found
					t.Fatalf(
						"Expected session (%s) not found in: %v",
						createdSession.Key,
						sessionStorage,
					)
					return nil, nil
				},
				// Define dummy hook to enable sessions on this server
				OnSessionClosed: func(_ *webwire.Client) error { return nil },
			},
		},
	)

	// Initialize client
	client := webwireClient.NewClient(addr, webwireClient.Options{})

	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	/*****************************************************************\
		Step 1 - Create session and disconnect
	\*****************************************************************/

	// Create a new session
	if _, err := client.Request(
		"login",
		webwire.Payload{Data: []byte("auth")},
	); err != nil {
		t.Fatalf("Auth request failed: %s", err)
	}

	tmp := client.Session()
	createdSession = &tmp

	// Disconnect client without closing the session
	client.Close()

	// Ensure the session isn't lost
	if client.Status() == webwireClient.StatConnected {
		t.Fatal("Client is expected to be disconnected")
	}
	if client.Session().Key == "" {
		t.Fatal("Session lost after disconnection")
	}

	/*****************************************************************\
		Step 2 - Reconnect, restore and verify authentication
	\*****************************************************************/
	currentStep = 2

	// Reconnect (this should automatically try to restore the session)
	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't reconnect: %s", err)
	}

	// Verify whether the previous session was restored automatically
	// and the server authenticates the user
	if _, err := client.Request(
		"verify",
		webwire.Payload{Data: []byte("isrestored?")},
	); err != nil {
		t.Fatalf("Second request failed: %s", err)
	}
}
