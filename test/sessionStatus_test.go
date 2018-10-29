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
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Try to create a new session
				if err := conn.CreateSession(nil); err != nil {
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

	require.Equal(t, 0, server.ActiveSessionsNum())

	// Initialize client A
	clientA := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	// Authenticate and create session
	authReqReply, err := clientA.connection.Request(
		context.Background(),
		"login",
		wwr.NewPayload(wwr.EncodingBinary, []byte("bla")),
	)
	require.NoError(t, err)

	session := clientA.connection.Session()
	require.Equal(t, session.Key, string(authReqReply.Data()))

	// Check status, expect 1 session with 1 connection
	require.Equal(t, 1, server.ActiveSessionsNum())
	require.Equal(t, 1, server.SessionConnectionsNum(session.Key))
	require.Len(t, server.SessionConnections(session.Key), 1)

	// Initialize client B
	clientB := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	require.NoError(t, clientB.connection.RestoreSession(authReqReply.Data()))

	// Check status, expect 1 session with 2 connections
	require.Equal(t, 1, server.ActiveSessionsNum())
	require.Equal(t, 2, server.SessionConnectionsNum(session.Key))
	require.Len(t, server.SessionConnections(session.Key), 2)

	// Close first connection
	require.NoError(t, clientA.connection.CloseSession())

	// Check status, expect 1 session with 1 connection
	require.Equal(t, 1, server.ActiveSessionsNum())
	require.Equal(t, 1, server.SessionConnectionsNum(session.Key))
	require.Len(t, server.SessionConnections(session.Key), 1)

	// Close second connection
	require.NoError(t, clientB.connection.CloseSession())

	// Check status, expect 0 sessions
	require.Equal(t, 0, server.ActiveSessionsNum())
	require.Equal(t, -1, server.SessionConnectionsNum(session.Key))
	require.Nil(t, server.SessionConnections(session.Key))
}
