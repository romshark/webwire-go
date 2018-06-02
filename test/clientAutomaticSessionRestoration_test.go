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
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				clt *webwire.Client,
				msg *webwire.Message,
			) (webwire.Payload, error) {
				if currentStep == 2 {
					// Expect the session to have been automatically restored
					compareSessions(t, createdSession, clt.Session())
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
			SessionManager: &callbackPoweredSessionManager{
				// Saves the session
				SessionCreated: func(client *webwire.Client) error {
					sess := client.Session()
					sessionStorage[sess.Key] = sess
					return nil
				},
				// Finds session by key
				SessionLookup: func(key string) (
					webwire.SessionLookupResult,
					error,
				) {
					// Expect the key of the created session to be looked up
					if key != createdSession.Key {
						err := fmt.Errorf(
							"Expected and looked up session keys differ: %s | %s",
							createdSession.Key,
							key,
						)
						t.Fatalf("Session lookup mismatch: %s", err)
						return webwire.SessionLookupResult{}, err
					}

					if session, exists := sessionStorage[key]; exists {
						// Session found
						return webwire.SessionLookupResult{
							Creation:   session.Creation,
							LastLookup: session.LastLookup,
							Info: webwire.SessionInfoToVarMap(
								session.Info,
							),
						}, nil
					}

					// Expect the session to be found
					t.Fatalf(
						"Expected session (%s) not found in: %v",
						createdSession.Key,
						sessionStorage,
					)

					// Session not found
					return webwire.SessionLookupResult{},
						webwire.SessNotFoundErr{}
				},
			},
		},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{},
		callbackPoweredClientHooks{},
	)

	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	/*****************************************************************\
		Step 1 - Create session and disconnect
	\*****************************************************************/

	// Create a new session
	if _, err := client.connection.Request(
		"login",
		webwire.Payload{Data: []byte("auth")},
	); err != nil {
		t.Fatalf("Auth request failed: %s", err)
	}

	createdSession = client.connection.Session()

	// Disconnect client without closing the session
	client.connection.Close()

	// Ensure the session isn't lost
	if client.connection.Status() == webwireClient.StatConnected {
		t.Fatal("Client is expected to be disconnected")
	}
	if client.connection.Session().Key == "" {
		t.Fatal("Session lost after disconnection")
	}

	/*****************************************************************\
		Step 2 - Reconnect, restore and verify authentication
	\*****************************************************************/
	currentStep = 2

	// Reconnect (this should automatically try to restore the session)
	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't reconnect: %s", err)
	}

	// Verify whether the previous session was restored automatically
	// and the server authenticates the user
	if _, err := client.connection.Request(
		"verify",
		webwire.Payload{Data: []byte("isrestored?")},
	); err != nil {
		t.Fatalf("Second request failed: %s", err)
	}
}
