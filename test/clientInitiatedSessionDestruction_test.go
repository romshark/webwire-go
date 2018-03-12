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

// TestClientInitiatedSessionDestruction verifies
// client-initiated session destruction
func TestClientInitiatedSessionDestruction(t *testing.T) {
	sessionCreationCallbackCalled := NewPending(1, 1*time.Second, true)
	sessionDestructionCallbackCalled := NewPending(1, 1*time.Second, true)
	var createdSession *webwire.Session
	expectedCredentials := webwire.Payload{
		Encoding: webwire.EncodingUtf8,
		Data:     []byte("secret_credentials"),
	}
	placeholderMessage := webwire.Payload{
		Encoding: webwire.EncodingUtf8,
		Data:     []byte("nothinginteresting"),
	}
	currentStep := 1

	// Initialize webwire server
	_, addr := setupServer(
		t,
		webwire.Hooks{
			OnRequest: func(ctx context.Context) (webwire.Payload, error) {
				// Extract request message and requesting client from the context
				msg := ctx.Value(webwire.Msg).(webwire.Message)

				// On step 2 - verify session creation and correctness
				if currentStep == 2 {
					compareSessions(t, createdSession, msg.Client.Session)
					if string(msg.Payload.Data) != msg.Client.Session.Key {
						t.Errorf(
							"Clients session key doesn't match: "+
								"client: '%s' | server: '%s'",
							string(msg.Payload.Data),
							msg.Client.Session.Key,
						)
					}
					return webwire.Payload{}, nil
				}

				// On step 4 - verify session destruction
				if currentStep == 4 {
					if msg.Client.Session != nil {
						t.Errorf(
							"Expected the session to be destroyed, got: %v",
							msg.Client.Session,
						)
					}
					return webwire.Payload{}, nil
				}

				// On step 1 - authenticate and create a new session
				if err := msg.Client.CreateSession(nil); err != nil {
					return webwire.Payload{}, webwire.Error{
						Code:    "INTERNAL_ERROR",
						Message: fmt.Sprintf("Internal server error: %s", err),
					}
				}

				// Return the key of the newly created session
				return webwire.Payload{
					Data: []byte(msg.Client.Session.Key),
				}, nil
			},
			// Define dummy hooks to enable sessions on this server
			OnSessionCreated: func(_ *webwire.Client) error { return nil },
			OnSessionLookup:  func(_ string) (*webwire.Session, error) { return nil, nil },
			OnSessionClosed:  func(_ *webwire.Client) error { return nil },
		},
	)

	// Initialize client
	client := webwireClient.NewClient(
		addr,
		webwireClient.Hooks{
			OnSessionCreated: func(_ *webwire.Session) {
				// Mark the client-side session creation callback as executed
				sessionCreationCallbackCalled.Done()
			},
			OnSessionClosed: func() {
				// Ensure this callback is called during the
				if currentStep != 3 {
					t.Errorf(
						"Client-side session destruction callback "+
							"called at wrong step (%d)",
						currentStep,
					)
				}
				sessionDestructionCallbackCalled.Done()
			},
		},
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)

	/*****************************************************************\
		Step 1 - Session Creation
	\*****************************************************************/
	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send authentication request
	authReqReply, err := client.Request("login", expectedCredentials)
	if err != nil {
		t.Fatalf("Authentication request failed: %s", err)
	}

	tmp := client.Session()
	createdSession = &tmp

	// Verify reply
	if createdSession.Key != string(authReqReply.Data) {
		t.Fatalf(
			"Unexpected session key: %s | %s",
			createdSession.Key,
			string(authReqReply.Data),
		)
	}

	// Wait for the client-side session creation callback to be executed
	if err := sessionCreationCallbackCalled.Wait(); err != nil {
		t.Fatal("Session creation callback not called")
	}

	// Ensure the session was locally created
	currentSessionAfterCreation := client.Session()
	if currentSessionAfterCreation.Key == "" {
		t.Fatalf(
			"Expected session on client-side, got none: %v",
			currentSessionAfterCreation,
		)
	}

	/*****************************************************************\
		Step 2 - Session Creation Verification
	\*****************************************************************/
	currentStep = 2

	// Send a test-request to verify the session creation on the server
	if _, err := client.Request(
		"verify-session-created",
		webwire.Payload{Data: []byte(client.Session().Key)},
	); err != nil {
		t.Fatalf("Session creation verification request failed: %s", err)
	}

	/*****************************************************************\
		Step 3 - Client-Side Session Destruction
	\*****************************************************************/
	currentStep = 3

	// Request session destruction
	if err := client.CloseSession(); err != nil {
		t.Fatalf("Failed closing session on the client: %s", err)
	}

	// Wait for the client-side session destruction callback to be called
	if err := sessionDestructionCallbackCalled.Wait(); err != nil {
		t.Fatal("Session destruction callback not called")
	}

	/*****************************************************************\
		Step 4 - Destruction Verification
	\*****************************************************************/
	currentStep = 4

	// Ensure the session is destroyed locally as well
	currentSessionAfterDestruction := client.Session()
	if currentSessionAfterDestruction.Key != "" {
		t.Fatalf(
			"Expected session to be destroyed on the client as well, "+
				"but still got: %v",
			currentSessionAfterDestruction,
		)
	}

	// Send a test-request to verify the session was destroyed on the server
	if _, err := client.Request("test-request", placeholderMessage); err != nil {
		t.Fatalf("Session destruction verification request failed: %s", err)
	}
}
