package test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestAuthentication verifies the server is connectable,
// and is able to receives requests and signals, create sessions
// and identify clients during request- and signal handling
func TestAuthentication(t *testing.T) {
	// Because compareSessions doesn't compare the sessions attached info:
	compareSessionInfo := func(actual *wwr.Session) {
		info, ok := actual.Info.(struct {
			UserID     string `json:"uid"`
			SomeNumber int    `json:"some-number"`
		})
		if !ok {
			t.Errorf("Couldn't cast info from: %v", actual.Info)
		}

		// Check uid
		field := "session.info.UserID"
		expectedUID := "clientidentifiergoeshere"
		if info.UserID != expectedUID {
			t.Errorf("%s differs: %s | %s", field, info.UserID, expectedUID)
			return
		}

		// Check some-number
		field = "session.info.some-number"
		expectedNumber := int(12345)
		if info.SomeNumber != expectedNumber {
			t.Errorf("%s differs: %d | %d", field, info.SomeNumber, expectedNumber)
			return
		}
	}

	onSessionCreatedHookExecuted := NewPending(1, 1*time.Second, true)
	clientSignalReceived := NewPending(1, 1*time.Second, true)
	var createdSession *wwr.Session
	sessionInfo := struct {
		UserID     string `json:"uid"`
		SomeNumber int    `json:"some-number"`
	}{
		"clientidentifiergoeshere",
		12345,
	}
	expectedCredentials := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     []byte("secret_credentials"),
	}
	expectedConfirmation := wwr.Payload{
		Encoding: wwr.EncodingUtf8,
		Data:     []byte("session_is_correct"),
	}
	currentStep := 1

	// Initialize webwire server
	_, addr := setupServer(
		t,
		wwr.Hooks{
			OnSignal: func(ctx context.Context) {
				defer clientSignalReceived.Done()
				// Extract request message and requesting client from the context
				msg := ctx.Value(wwr.Msg).(wwr.Message)
				compareSessions(t, createdSession, msg.Client.Session)
				compareSessionInfo(msg.Client.Session)
			},
			OnRequest: func(ctx context.Context) (wwr.Payload, *wwr.Error) {
				// Extract request message and requesting client from the context
				msg := ctx.Value(wwr.Msg).(wwr.Message)

				// If already authenticated then check session
				if currentStep > 1 {
					compareSessions(t, createdSession, msg.Client.Session)
					compareSessionInfo(msg.Client.Session)
					return expectedConfirmation, nil
				}

				// Try to create a new session
				if err := msg.Client.CreateSession(sessionInfo); err != nil {
					return wwr.Payload{}, &wwr.Error{
						Code:    "INTERNAL_ERROR",
						Message: fmt.Sprintf("Internal server error: %s", err),
					}
				}

				// Authentication step is passed
				currentStep = 2

				// Return the key of the newly created session (use default binary encoding)
				return wwr.Payload{
					Data: []byte(msg.Client.Session.Key),
				}, nil
			},
			// Define dummy hooks to enable sessions on this server
			OnSessionCreated: func(_ *wwr.Client) error { return nil },
			OnSessionLookup:  func(_ string) (*wwr.Session, error) { return nil, nil },
			OnSessionClosed:  func(_ *wwr.Client) error { return nil },
		},
	)

	// Initialize client
	client := wwrclt.NewClient(
		addr,
		wwrclt.Hooks{
			OnSessionCreated: func(session *wwr.Session) {
				// The session info object won't be of initial structure type
				// because of intermediate JSON encoding
				// it'll be a map of arbitrary values with string keys
				info := session.Info.(map[string]interface{})

				// Check uid
				field := "session.info.uid"
				expectedUID := "clientidentifiergoeshere"
				actualUID, ok := info["uid"].(string)
				if !ok {
					t.Errorf("expected %s not string", field)
					return
				}
				if actualUID != expectedUID {
					t.Errorf("%s differs: %s | %s", field, actualUID, expectedUID)
					return
				}

				// Check some-number
				field = "session.info.some-number"
				expectedNumber := float64(12345)
				actualNumber, ok := info["some-number"].(float64)
				if !ok {
					t.Errorf("expected %s not float64", field)
					return
				}
				if actualNumber != expectedNumber {
					t.Errorf("%s differs: %f | %f", field, actualNumber, expectedNumber)
					return
				}
				onSessionCreatedHookExecuted.Done()
			},
		},
		5*time.Second,
		os.Stdout,
		os.Stderr,
	)
	defer client.Close()

	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send authentication request and await reply
	authReqReply, err := client.Request("login", expectedCredentials)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	tmp := client.Session()
	createdSession = &tmp

	// Verify reply
	comparePayload(
		t,
		"authentication reply",
		wwr.Payload{
			Data: []byte(createdSession.Key),
		},
		authReqReply,
	)

	// Send a test-request to verify the session on the server
	// and await response
	testReqReply, err := client.Request("test", expectedCredentials)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify reply
	comparePayload(t, "test reply", expectedConfirmation, testReqReply)

	// Send a test-signal to verify the session on the server
	if err := client.Signal("test", expectedCredentials); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	if err := clientSignalReceived.Wait(); err != nil {
		t.Fatal("Client signal not received")
	}

	// Expect the session creation hook to be executed in the client
	if err := onSessionCreatedHookExecuted.Wait(); err != nil {
		t.Fatalf("client.OnSessionCreated hook wasn't executed")
	}
}
