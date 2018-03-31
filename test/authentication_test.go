package test

import (
	"context"
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestAuthentication tests session creation and client authentication
// during request- and signal handling
func TestAuthentication(t *testing.T) {
	// Because compareSessions doesn't compare the sessions attached info:
	compareSessionInfo := func(actual *wwr.Session) {
		// Check uid
		field := "session.info.UserID"
		expectedUserIdent := "clientidentifiergoeshere"
		actualUserIdent, correctType := actual.Info["uid"].(string)
		if !correctType {
			t.Errorf("%s incorrect type: %s", field, reflect.TypeOf(actual.Info["uid"]))
		}
		if actualUserIdent != expectedUserIdent {
			t.Errorf("%s differs: %s | %s", field, actualUserIdent, expectedUserIdent)
			return
		}

		// Check some-number
		field = "session.info.some-number"
		expectedNumber := int(12345)
		actualNumber, correctType := actual.Info["some-number"].(int)
		if !correctType {
			t.Errorf("%s incorrect type: %s", field, reflect.TypeOf(actual.Info["some-number"]))
		}
		if actualNumber != expectedNumber {
			t.Errorf("%s differs: %d | %d", field, actualNumber, expectedNumber)
			return
		}
	}

	onSessionCreatedHookExecuted := NewPending(1, 1*time.Second, true)
	clientSignalReceived := NewPending(1, 1*time.Second, true)
	var createdSession *wwr.Session
	sessionInfo := make(wwr.SessionInfo)
	sessionInfo["uid"] = "clientidentifiergoeshere"
	sessionInfo["some-number"] = 12345
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
	server := setupServer(
		t,
		wwr.ServerOptions{
			SessionsEnabled: true,
			Hooks: wwr.Hooks{
				OnSignal: func(ctx context.Context) {
					defer clientSignalReceived.Done()
					// Extract request message and requesting client from the context
					msg := ctx.Value(wwr.Msg).(wwr.Message)
					sess := msg.Client.Session()
					compareSessions(t, createdSession, sess)
					compareSessionInfo(sess)
				},
				OnRequest: func(ctx context.Context) (wwr.Payload, error) {
					// Extract request message and requesting client from the context
					msg := ctx.Value(wwr.Msg).(wwr.Message)

					// If already authenticated then check session
					if currentStep > 1 {
						sess := msg.Client.Session()
						compareSessions(t, createdSession, sess)
						compareSessionInfo(sess)
						return expectedConfirmation, nil
					}

					// Try to create a new session
					if err := msg.Client.CreateSession(sessionInfo); err != nil {
						return wwr.Payload{}, err
					}

					// Authentication step is passed
					currentStep = 2

					// Return the key of the newly created session (use default binary encoding)
					return wwr.Payload{
						Data: []byte(msg.Client.SessionKey()),
					}, nil
				},
			},
		},
	)

	// Initialize client
	client := wwrclt.NewClient(
		server.Addr().String(),
		wwrclt.Options{
			Hooks: wwrclt.Hooks{
				OnSessionCreated: func(session *wwr.Session) {
					// The session info object won't be of initial structure type
					// because of intermediate JSON encoding
					// it'll be a map of arbitrary values with string keys
					info := session.Info

					// Check uid
					field := "session.info.uid"
					expectedUserIdent := "clientidentifiergoeshere"
					actualUID, ok := info["uid"].(string)
					if !ok {
						t.Errorf("expected %s not string", field)
						return
					}
					if actualUID != expectedUserIdent {
						t.Errorf("%s differs: %s | %s", field, actualUID, expectedUserIdent)
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
			DefaultRequestTimeout: 2 * time.Second,
		},
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
