package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
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
				conn wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				// Close session on logout
				if string(msg.Name()) == "logout" {
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
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: time.Second * 2,
		},
		callbackPoweredClientHooks{},
	)
	defer client.connection.Close()

	require.NoError(t, client.connection.Connect())

	// Send authentication request
	_, err := client.connection.Request(
		context.Background(),
		[]byte("login"),
		wwr.NewPayload(wwr.EncodingUtf8, []byte("nothing")),
	)
	require.NoError(t, err)

	activeSessionNumberBefore := server.ActiveSessionsNum()
	require.Equal(t,
		1, activeSessionNumberBefore,
		"Unexpected active session number after authentication",
	)

	// Send logout request
	_, err = client.connection.Request(
		context.Background(),
		[]byte("logout"),
		wwr.NewPayload(wwr.EncodingUtf8, []byte("nothing")),
	)
	require.NoError(t, err)

	activeSessionNumberAfter := server.ActiveSessionsNum()
	require.Equal(t,
		0, activeSessionNumberAfter,
		"Unexpected active session number after logout",
	)
}
