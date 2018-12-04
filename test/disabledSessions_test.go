package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDisabledSessions tests errors returned by CreateSession, CloseSession and
// client.RestoreSession when sessions are disabled
func TestDisabledSessions(t *testing.T) {
	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Try to create a new session and expect an error
				createErr := conn.CreateSession(nil)
				assert.IsType(t, wwr.SessionsDisabledErr{}, createErr)

				// Try to create a new session and expect an error
				closeErr := conn.CloseSession()
				assert.IsType(t, wwr.SessionsDisabledErr{}, closeErr)

				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{
			Sessions: wwr.Disabled,
		},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		TestClientHooks{
			OnSessionCreated: func(*wwr.Session) {
				t.Errorf("OnSessionCreated was not expected to be called")
			},
		},
	)

	// Send authentication request and await reply
	_, err := client.Connection.Request(
		context.Background(),
		[]byte("login"),
		wwr.Payload{Data: []byte("testdata")},
	)
	require.NoError(t, err)

	sessRestErr := client.Connection.RestoreSession(
		context.Background(),
		[]byte("testkey"),
	)
	assert.IsType(t, wwr.SessionsDisabledErr{}, sessRestErr)
}
