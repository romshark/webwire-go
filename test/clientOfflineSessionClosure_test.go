package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientOfflineSessionClosure tests offline session closure
func TestClientOfflineSessionClosure(t *testing.T) {
	sessionStorage := make(map[string]*webwire.Session)

	currentStep := 1
	var createdSession *webwire.Session

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				clt *webwire.Client,
				_ *webwire.Message,
			) (webwire.Payload, error) {
				if currentStep == 2 {
					// Expect the session to be removed
					if clt.HasSession() {
						t.Errorf("Expected client to be anonymous")
					}
					return webwire.Payload{}, nil
				}

				// Try to create a new session
				if err := clt.CreateSession(nil); err != nil {
					return webwire.Payload{}, err
				}

				// Return the key of the newly created session
				return webwire.Payload{}, nil
			},
		},
		webwire.ServerOptions{
			SessionsEnabled: true,
			SessionManager: &CallbackPoweredSessionManager{
				// Saves the session
				SessionCreated: func(client *webwire.Client) error {
					sess := client.Session()
					sessionStorage[sess.Key] = sess
					return nil
				},
				// Finds session by key
				SessionLookup: func(key string) (*webwire.Session, error) {
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
			},
		},
	)

	// Initialize client
	client := webwireClient.NewClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)

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
		Step 2 - Close session, reconnect and verify
	\*****************************************************************/
	currentStep = 2

	if err := client.CloseSession(); err != nil {
		t.Fatalf("Offline session closure failed: %s", err)
	}

	// Ensure the session is removed locally
	if client.Session().Key != "" {
		t.Fatal("Session not removed")
	}

	// Reconnect
	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't reconnect: %s", err)
	}

	// Ensure the client is anonymous
	if _, err := client.Request(
		"verify-restored",
		webwire.Payload{Data: []byte("isrestored?")},
	); err != nil {
		t.Fatalf("Second request failed: %s", err)
	}
}
