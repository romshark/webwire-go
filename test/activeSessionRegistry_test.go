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

// TestActiveSessionRegistry verifies that the session registry
// of currently active sessions is properly updated
func TestActiveSessionRegistry(t *testing.T) {
	// Initialize webwire server
	setup := setupTestServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				// Close session on logout
				if string(msg.Name()) == "logout" {
					assert.NoError(t, conn.CloseSession())
					return wwr.Payload{}, nil
				}

				// Try to create a new session
				err := conn.CreateSession(nil)
				assert.NoError(t, err)
				if err != nil {
					return wwr.Payload{}, err
				}

				// Return the key of the newly created session
				// (use default binary encoding)
				return wwr.Payload{
					Encoding: wwr.EncodingBinary,
					Data:     []byte(conn.SessionKey()),
				}, nil
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.newClient(
		wwrclt.Options{
			DefaultRequestTimeout: time.Second * 2,
		},
		nil, // Use the default transport implementation
		testClientHooks{},
	)
	defer client.connection.Close()

	require.NoError(t, client.connection.Connect())

	// Send authentication request
	reply, err := client.connection.Request(
		context.Background(),
		[]byte("login"),
		wwr.Payload{
			Encoding: wwr.EncodingUtf8,
			Data:     []byte("nothing"),
		},
	)
	require.NoError(t, err)
	reply.Close()

	activeSessionNumberBefore := setup.Server.ActiveSessionsNum()
	require.Equal(t,
		1, activeSessionNumberBefore,
		"Unexpected active session number after authentication",
	)

	// Send logout request
	reply, err = client.connection.Request(
		context.Background(),
		[]byte("logout"),
		wwr.Payload{
			Encoding: wwr.EncodingUtf8,
			Data:     []byte("nothing"),
		},
	)
	require.NoError(t, err)
	reply.Close()

	activeSessionNumberAfter := setup.Server.ActiveSessionsNum()
	require.Equal(t,
		0, activeSessionNumberAfter,
		"Unexpected active session number after logout",
	)
}
