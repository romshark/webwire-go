package test

import (
	"context"
	"sync"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testAuthenticationSessInfo struct {
	UserIdent  string
	SomeNumber int
}

// Copy implements the webwire.SessionInfo interface.
// It deep-copies the object and returns it's exact clone
func (sinf *testAuthenticationSessInfo) Copy() wwr.SessionInfo {
	return &testAuthenticationSessInfo{
		UserIdent:  sinf.UserIdent,
		SomeNumber: sinf.SomeNumber,
	}
}

// Fields implements the webwire.SessionInfo interface.
// It returns a constant list of the names of all fields of the object
func (sinf *testAuthenticationSessInfo) Fields() []string {
	return []string{"uid", "some-number"}
}

// Copy implements the webwire.SessionInfo interface.
// It deep-copies the field identified by the provided name
// and returns it's exact clone
func (sinf *testAuthenticationSessInfo) Value(fieldName string) interface{} {
	switch fieldName {
	case "uid":
		return sinf.UserIdent
	case "some-number":
		return sinf.SomeNumber
	}
	return nil
}

// TestAuthentication tests session creation and client authentication during
// request- and signal handling
func TestAuthentication(t *testing.T) {
	// Because CompareSessions doesn't compare the sessions attached info:
	compareSessionInfo := func(actual *wwr.Session) {
		// Check uid
		expectedUserIdent := "clientidentifiergoeshere"

		assert.IsType(t, "string", actual.Info.Value("uid"))
		actualUserIdent := actual.Info.Value("uid").(string)

		assert.Equal(t, expectedUserIdent, actualUserIdent)

		// Check some-number
		expectedNumber := int(12345)

		assert.IsType(t, expectedNumber, actual.Info.Value("some-number"))
		actualNumber := actual.Info.Value("some-number").(int)
		assert.Equal(t, expectedNumber, actualNumber)
	}

	onSessionCreatedHookExecuted := sync.WaitGroup{}
	onSessionCreatedHookExecuted.Add(1)
	clientSignalReceived := sync.WaitGroup{}
	clientSignalReceived.Add(1)
	var createdSession *wwr.Session
	sessionInfo := &testAuthenticationSessInfo{
		UserIdent:  "clientidentifiergoeshere",
		SomeNumber: 12345,
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
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Signal: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) {
				defer clientSignalReceived.Done()
				sess := conn.Session()
				CompareSessions(t, createdSession, sess)
				compareSessionInfo(sess)
			},
			Request: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// If already authenticated then check session
				if currentStep > 1 {
					sess := conn.Session()
					CompareSessions(t, createdSession, sess)
					compareSessionInfo(sess)
					return expectedConfirmation, nil
				}

				// Try to create a new session
				err := conn.CreateSession(sessionInfo)
				assert.NoError(t, err)
				if err != nil {
					return wwr.Payload{}, err
				}

				// Authentication step is passed
				currentStep = 2

				// Return the key of the newly created session
				// (use default binary encoding)
				return wwr.Payload{Data: []byte(conn.SessionKey())}, nil
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
			SessionInfoParser: func(
				data map[string]interface{},
			) wwr.SessionInfo {
				return &testAuthenticationSessInfo{
					UserIdent:  data["uid"].(string),
					SomeNumber: int(data["some-number"].(float64)),
				}
			},
		},
		nil, // Use the default transport implementation
		TestClientHooks{
			OnSessionCreated: func(session *wwr.Session) {
				// The session info object won't be of initial structure type
				// because of intermediate JSON encoding
				// it'll be a map of arbitrary values with string keys
				compareSessionInfo(session)
				onSessionCreatedHookExecuted.Done()
			},
		},
	)

	// Send authentication request and await reply
	authReqReply, err := client.Connection.Request(
		context.Background(),
		[]byte("login"),
		expectedCredentials,
	)
	require.NoError(t, err)

	createdSession = client.Connection.Session()

	// Verify reply
	require.Equal(t, wwr.EncodingBinary, authReqReply.PayloadEncoding())
	require.Equal(t, []byte(createdSession.Key), authReqReply.Payload())
	authReqReply.Close()

	// Send a test-request to verify the session on the server
	// and await response
	testReqReply, err := client.Connection.Request(
		context.Background(),
		[]byte("test"),
		expectedCredentials,
	)
	require.NoError(t, err)

	// Verify reply
	require.Equal(
		t,
		expectedConfirmation.Encoding,
		testReqReply.PayloadEncoding(),
	)
	require.Equal(
		t,
		expectedConfirmation.Data,
		testReqReply.Payload(),
	)

	// Send a test-signal to verify the session on the server
	require.NoError(t, client.Connection.Signal(
		context.Background(),
		[]byte("test"),
		expectedCredentials,
	))

	clientSignalReceived.Wait()

	// Expect the session creation hook to be executed in the client
	onSessionCreatedHookExecuted.Wait()
}
