package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestSessionStatus tests session monitoring methods
func TestSessionStatus(t *testing.T) {
	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Try to create a new session
				if err := conn.CreateSession(nil); err != nil {
					return wwr.Payload{}, err
				}

				// Return the key of the newly created session
				// (use default binary encoding)
				return wwr.Payload{Data: []byte(conn.SessionKey())}, nil
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	require.Equal(t, 0, setup.Server.ActiveSessionsNum())

	// Initialize client A
	clientA := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		TestClientHooks{},
	)

	// Authenticate and create session
	authReqReply, err := clientA.Connection.Request(
		context.Background(),
		[]byte("login"),
		wwr.Payload{Data: []byte("bla")},
	)
	require.NoError(t, err)

	session := clientA.Connection.Session()
	require.Equal(t, session.Key, string(authReqReply.Payload()))

	// Check status, expect 1 session with 1 connection
	require.Equal(t, 1, setup.Server.ActiveSessionsNum())
	require.Equal(t, 1, setup.Server.SessionConnectionsNum(session.Key))
	require.Len(t, setup.Server.SessionConnections(session.Key), 1)

	// Initialize client B
	clientB := setup.NewClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		TestClientHooks{},
	)

	require.NoError(t, clientB.Connection.RestoreSession(
		context.Background(),
		authReqReply.Payload(),
	))

	// Check status, expect 1 session with 2 connections
	require.Equal(t, 1, setup.Server.ActiveSessionsNum())
	require.Equal(t, 2, setup.Server.SessionConnectionsNum(session.Key))
	require.Len(t, setup.Server.SessionConnections(session.Key), 2)

	// Close first connection
	require.NoError(t, clientA.Connection.CloseSession())

	// Check status, expect 1 session with 1 connection
	require.Equal(t, 1, setup.Server.ActiveSessionsNum())
	require.Equal(t, 1, setup.Server.SessionConnectionsNum(session.Key))
	require.Len(t, setup.Server.SessionConnections(session.Key), 1)

	// Close second connection
	require.NoError(t, clientB.Connection.CloseSession())

	// Check status, expect 0 sessions
	require.Equal(t, 0, setup.Server.ActiveSessionsNum())
	require.Equal(t, -1, setup.Server.SessionConnectionsNum(session.Key))
	require.Nil(t, setup.Server.SessionConnections(session.Key))
}
