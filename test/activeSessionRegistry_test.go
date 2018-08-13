package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestActiveSessionRegistry verifies that the session registry
// of currently active sessions is properly updated
func TestActiveSessionRegistry(t *testing.T) {
	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn webwire.Connection,
				msg webwire.Message,
			) (webwire.Payload, error) {
				// Close session on logout
				if msg.Name() == "logout" {
					assert.NoError(t, conn.CloseSession())
					return nil, nil
				}

				// Try to create a new session
				err := conn.CreateSession(nil)
				assert.NoError(t, err)
				if err != nil {
					return nil, err
				}

				// Return the key of the newly created session
				// (use default binary encoding)
				return webwire.NewPayload(
					webwire.EncodingBinary,
					[]byte(conn.SessionKey()),
				), nil
			},
		},
		webwire.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: time.Second * 2,
		},
		callbackPoweredClientHooks{},
	)
	defer client.connection.Close()

	require.NoError(t, client.connection.Connect())

	// Send authentication request
	_, err := client.connection.Request(
		context.Background(),
		"login",
		webwire.NewPayload(webwire.EncodingUtf8, []byte("nothing")),
	)
	require.NoError(t, err)

	activeSessionNumberBefore := server.ActiveSessionsNum()
	require.Equal(t,
		1, activeSessionNumberBefore,
		"Unexpected active session number after authentication: %d",
		activeSessionNumberBefore,
	)

	// Send logout request
	_, err = client.connection.Request(
		context.Background(),
		"logout",
		webwire.NewPayload(webwire.EncodingUtf8, []byte("nothing")),
	)
	require.NoError(t, err)

	activeSessionNumberAfter := server.ActiveSessionsNum()
	require.Equal(t,
		0, activeSessionNumberAfter,
		"Unexpected active session number after logout: %d",
		activeSessionNumberAfter,
	)
}
