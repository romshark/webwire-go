package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientSessionRestoration tests manual session restoration by key
func TestClientSessionRestoration(t *testing.T) {
	sessionStorage := make(map[string]*webwire.Session)

	currentStep := 1
	var createdSession *webwire.Session

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn webwire.Connection,
				_ webwire.Message,
			) (webwire.Payload, error) {
				if currentStep == 2 {
					// Expect the session to be automatically restored
					compareSessions(t, createdSession, conn.Session())
					return nil, nil
				}

				// Try to create a new session
				if err := conn.CreateSession(nil); err != nil {
					return nil, err
				}

				// Return the key of the newly created session
				return nil, nil
			},
		},
		webwire.ServerOptions{
			SessionManager: &callbackPoweredSessionManager{
				// Saves the session
				SessionCreated: func(conn webwire.Connection) error {
					sess := conn.Session()
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
							"Expected and found session keys differ: %s | %s",
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

					//Session not found
					return webwire.SessionLookupResult{},
						webwire.SessNotFoundErr{}
				},
			},
		},
	)

	// Initialize client
	initialClient := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	if err := initialClient.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect initial client: %s", err)
	}

	/*****************************************************************\
		Step 1 - Create session and disconnect
	\*****************************************************************/

	// Create a new session
	if _, err := initialClient.connection.Request(
		context.Background(),
		"login",
		webwire.NewPayload(webwire.EncodingBinary, []byte("auth")),
	); err != nil {
		t.Fatalf("Auth request failed: %s", err)
	}

	createdSession = initialClient.connection.Session()

	// Disconnect client without closing the session
	initialClient.connection.Close()

	/*****************************************************************\
		Step 2 - Create new client, restore session from key
	\*****************************************************************/
	currentStep = 2

	// Initialize client
	secondClient := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	if err := secondClient.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect second client: %s", err)
	}

	// Ensure there's no active session on the second client
	sessionBefore := secondClient.connection.Session()
	if sessionBefore != nil {
		t.Fatalf(
			"Expected the second client to have no session, got: %v",
			sessionBefore,
		)
	}

	// Try to manually restore the session
	// using the initial clients session key
	if err := secondClient.connection.RestoreSession(
		[]byte(createdSession.Key),
	); err != nil {
		t.Fatalf("Manual session restoration failed: %s", err)
	}

	// Verify session
	sessionAfter := secondClient.connection.Session()
	if sessionAfter.Key == "" {
		t.Fatalf(
			"Expected the second client to have an active session, got: %v",
			sessionAfter,
		)
	}
	compareSessions(t, createdSession, sessionAfter)
}
