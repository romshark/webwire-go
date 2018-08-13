package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
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

// TestAuthentication tests session creation and client authentication
// during request- and signal handling
func TestAuthentication(t *testing.T) {
	// Because compareSessions doesn't compare the sessions attached info:
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

	onSessionCreatedHookExecuted := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	clientSignalReceived := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	var createdSession *wwr.Session
	sessionInfo := &testAuthenticationSessInfo{
		UserIdent:  "clientidentifiergoeshere",
		SomeNumber: 12345,
	}
	expectedCredentials := wwr.NewPayload(
		wwr.EncodingUtf8,
		[]byte("secret_credentials"),
	)
	expectedConfirmation := wwr.NewPayload(
		wwr.EncodingUtf8,
		[]byte("session_is_correct"),
	)
	currentStep := 1

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onSignal: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) {
				defer clientSignalReceived.Progress(1)
				sess := conn.Session()
				compareSessions(t, createdSession, sess)
				compareSessionInfo(sess)
			},
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// If already authenticated then check session
				if currentStep > 1 {
					sess := conn.Session()
					compareSessions(t, createdSession, sess)
					compareSessionInfo(sess)
					return expectedConfirmation, nil
				}

				// Try to create a new session
				err := conn.CreateSession(sessionInfo)
				assert.NoError(t, err)
				if err != nil {
					return nil, err
				}

				// Authentication step is passed
				currentStep = 2

				// Return the key of the newly created session
				// (use default binary encoding)
				return wwr.NewPayload(
					wwr.EncodingBinary,
					[]byte(conn.SessionKey()),
				), nil
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
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
		callbackPoweredClientHooks{
			OnSessionCreated: func(session *wwr.Session) {
				// The session info object won't be of initial structure type
				// because of intermediate JSON encoding
				// it'll be a map of arbitrary values with string keys
				compareSessionInfo(session)
				onSessionCreatedHookExecuted.Progress(1)
			},
		},
	)
	defer client.connection.Close()

	require.NoError(t, client.connection.Connect())

	// Send authentication request and await reply
	authReqReply, err := client.connection.Request(
		context.Background(),
		"login",
		expectedCredentials,
	)
	require.NoError(t, err)

	createdSession = client.connection.Session()

	// Verify reply
	comparePayload(t,
		wwr.NewPayload(
			wwr.EncodingBinary,
			[]byte(createdSession.Key),
		),
		authReqReply,
	)

	// Send a test-request to verify the session on the server
	// and await response
	testReqReply, err := client.connection.Request(
		context.Background(),
		"test",
		expectedCredentials,
	)
	require.NoError(t, err)

	// Verify reply
	comparePayload(t, expectedConfirmation, testReqReply)

	// Send a test-signal to verify the session on the server
	require.NoError(t, client.connection.Signal(
		"test",
		expectedCredentials,
	))

	require.NoError(t,
		clientSignalReceived.Wait(),
		"Client signal not received",
	)

	// Expect the session creation hook to be executed in the client
	require.NoError(t,
		onSessionCreatedHookExecuted.Wait(),
		"client.OnSessionCreated hook wasn't executed",
	)
}
