package test

import (
	"testing"
	"os"
	"fmt"
	"time"
	"sync"
	"context"

	webwire "github.com/qbeon/webwire-go"
	webwire_client "github.com/qbeon/webwire-go/client"
)

// TestClientInitiatedSessionDestruction verifies client-initiated session destruction
func TestClientInitiatedSessionDestruction(t *testing.T) {
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
		// onRequest
		func(ctx context.Context) ([]byte, *webwire.Error) {
			// Extract request message and requesting client from the context
			msg := ctx.Value(webwire.MESSAGE).(webwire.Message)

			// On step 2 - verify session creation and correctness
			if currentStep == 2 {
				compareSessions(t, createdSession, msg.Client.Session)
				if string(msg.Payload) != msg.Client.Session.Key {
					t.Errorf(
						"Clients session key doesn't match servers session key:" +
							" client: '%s' | server: '%s'",
						string(msg.Payload),
						msg.Client.Session.Key,
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
				webwire.Os_UNKNOWN,
				"user agent",
				nil,
			)
			createdSession = &newSession

			// Try to register the newly created session and bind it to the client
			if err := msg.Client.CreateSession(createdSession); err != nil {
				return nil, &webwire.Error {
					"INTERNAL_ERROR",
					fmt.Sprintf("Internal server error: %s", err),
				}
			}

			// Return the key of the newly created session
			return []byte(createdSession.Key), nil
		},
		// OnSaveSession
		func(session *webwire.Session) error {
			// Verify the session
			compareSessions(t, createdSession, session)
			return nil
		},
		// OnFindSession
		func(_ string) (*webwire.Session, error) {
			return nil, nil
		},
		// OnSessionClosure
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
		func (_ *webwire.Session) {
			// Mark the client-side session creation callback as executed
			sessionCreationCallback.Done()
		},
		// On session closed
		func() {
			// Ensure this callback is called during the 
			if currentStep != 3 {
				t.Errorf(
					"Client-side session destruction callback called at wrong step (%d)",
					currentStep,
				)
			}
			sessionDestructionCallback.Done()
		},
		5 * time.Second,
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
	comparePayload(t, "authentication reply", []byte(createdSession.Key), authReqReply)

	// Wait for the client-side session creation callback to be executed
	sessionCreationCallback.Wait()
	
	// Ensure the session was locally created
	currSessAfterCreation := client.Session()
	if currSessAfterCreation.Key == "" {
		t.Fatalf("Expected session on client-side, got none: %v", currSessAfterCreation)
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
		Step 3 - Client-Side Session Destruction
	\*****************************************************************/
	currentStep = 3

	// Request session destruction
	if err := client.CloseSession(); err != nil {
		t.Fatalf("Failed closing session on the client: %s", err)
	}

	// Wait for the client-side session destruction callback to be called
	sessionDestructionCallback.Wait()

	/*****************************************************************\
		Step 4 - Destruction Verification
	\*****************************************************************/
	currentStep = 4

	// Ensure the session is destroyed locally as well
	currSessAfterDestruction := client.Session()
	if currSessAfterDestruction.Key != "" {
		t.Fatalf(
			"Expected session to be destroyed on the client as well, but still got: %v",
			currSessAfterDestruction,
		)
	}

	// Send a test-request to verify the session was destroyed on the server
	if _, err := client.Request(placeholderMessage); err != nil {
		t.Fatalf("Session destruction verification request failed: %s", err)
	}
}