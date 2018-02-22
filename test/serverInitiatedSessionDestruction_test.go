package test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwire_client "github.com/qbeon/webwire-go/client"
	"github.com/qbeon/webwire-go/ostype"
)

// TestServerInitiatedSessionDestruction verifies
// server-initiated session destruction
func TestServerInitiatedSessionDestruction(t *testing.T) {
	var sessionCreationCallback sync.WaitGroup
	var sessionDestructionCallback sync.WaitGroup
	sessionCreationCallback.Add(1)
	sessionDestructionCallback.Add(1)
	var createdSession *webwire.Session
	expectedCredentials := []byte("secret_credentials")
	placeholderMessage := []byte("nothinginteresting")
	currentStep := 1

	// Initialize webwire server
	server := setupServer(
		t,
		nil,
		nil,
		nil,
		// onRequest
		func(ctx context.Context) ([]byte, *webwire.Error) {
			// Extract request message and requesting client from the context
			msg := ctx.Value(webwire.MESSAGE).(webwire.Message)

			// On step 2 - verify session creation and correctness
			if currentStep == 2 {
				compareSessions(t, createdSession, msg.Client.Session)
				if string(msg.Payload) != msg.Client.Session.Key {
					t.Errorf(
						"Clients session key doesn't match: "+
							"client: '%s' | server: '%s'",
						string(msg.Payload),
						msg.Client.Session.Key,
					)
				}
				return nil, nil
			}

			// on step 3 - close session and verify its destruction
			if currentStep == 3 {
				/***********************************************************\
					Server-side session destruction initiation
				\***********************************************************/
				// Attempt to destroy this clients session
				// on the end of the first step
				err := msg.Client.CloseSession()
				if err != nil {
					t.Errorf(
						"Couldn't close the active session on the server: %s",
						err,
					)
				}

				// Verify destruction
				if msg.Client.Session != nil {
					t.Errorf(
						"Expected the session to be destroyed, got: %v",
						msg.Client.Session,
					)
				}

				return nil, nil
			}

			// On step 4 - verify session destruction
			if currentStep == 4 {
				if msg.Client.Session != nil {
					t.Errorf("Expected the session to be destroyed")
				}
				return nil, nil
			}

			// On step 1 - authenticate and create a new session
			newSession := webwire.NewSession(
				ostype.Unknown,
				"user agent",
				nil,
			)
			createdSession = &newSession

			// Try to register the newly created session
			// and bind it to the client
			if err := msg.Client.CreateSession(createdSession); err != nil {
				return nil, &webwire.Error{
					Code:    "INTERNAL_ERROR",
					Message: fmt.Sprintf("Internal server error: %s", err),
				}
			}

			// Return the key of the newly created session
			return []byte(createdSession.Key), nil
		},
		// OnSessionCreated
		func(client *webwire.Client) error {
			// Verify the session
			compareSessions(t, createdSession, client.Session)
			return nil
		},
		// OnSessionLookup
		func(_ string) (*webwire.Session, error) {
			return nil, nil
		},
		// OnSessionClosed
		func(_ *webwire.Client) error {
			return nil
		},
		nil,
	)
	go server.Run()

	// Initialize client
	client := webwire_client.NewClient(
		server.Addr,
		nil,
		// On session created
		func(_ *webwire.Session) {
			// Mark the client-side session creation callback as executed
			sessionCreationCallback.Done()
		},
		// On session closed
		func() {
			// Ensure this callback is called during the
			if currentStep != 3 {
				t.Errorf(
					"Client-side session destruction callback "+
						"called at wrong step (%d)",
					currentStep,
				)
			}
			sessionDestructionCallback.Done()
		},
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)

	/*****************************************************************\
		Step 1 - Session Creation
	\*****************************************************************/
	// Send authentication request
	authReqReply, err := client.Request(expectedCredentials)
	if err != nil {
		t.Fatalf("Authentication request failed: %s", err)
	}

	// Verify reply
	comparePayload(
		t,
		"authentication reply",
		[]byte(createdSession.Key),
		authReqReply,
	)

	// Wait for the client-side session creation callback to be executed
	sessionCreationCallback.Wait()

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
	if _, err := client.Request([]byte(client.Session().Key)); err != nil {
		t.Fatalf("Session creation verification request failed: %s", err)
	}

	/*****************************************************************\
		Step 3 - Server-Side Session Destruction
	\*****************************************************************/
	currentStep = 3

	// Request session destruction
	if _, err := client.Request(placeholderMessage); err != nil {
		t.Fatalf("Session destruction request failed: %s", err)
	}

	// Wait for the client-side session destruction callback to be called
	sessionDestructionCallback.Wait()

	/*****************************************************************\
		Step 4 - Destruction Verification
	\*****************************************************************/
	currentStep = 4

	// Ensure the session is destroyed locally as well
	currentSessionAfterDestruction := client.Session()
	if currentSessionAfterDestruction.Key != "" {
		t.Fatalf(
			"Expected session to be destroyed on the client as well, got: %v",
			currentSessionAfterDestruction,
		)
	}

	// Send a test-request to verify the session was destroyed on the server
	if _, err := client.Request(placeholderMessage); err != nil {
		t.Fatalf("Session destruction verification request failed: %s", err)
	}
}
