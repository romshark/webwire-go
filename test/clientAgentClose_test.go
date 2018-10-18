package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestConnectionClose tests closing a client connection on the server-side
func TestConnectionClose(t *testing.T) {
	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				switch msg.Name() {
				case "closeA":
					fallthrough
				case "closeB":
					conn.Close()
					return nil, nil
				case "login":
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
				}
				return nil, fmt.Errorf("Invalid request %s", msg.Name())
			},
		},
		wwr.ServerOptions{},
	)

	actSess := server.ActiveSessionsNum()
	require.Equal(t, 0, actSess, "Unexpected number of active sessions")

	// Initialize client A
	clientA := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
			// Disable autoconnect to prevent automatic reconnection
			Autoconnect: wwr.Disabled,
		},
		callbackPoweredClientHooks{},
		nil, // No TLS configuration
	)

	require.NoError(t, clientA.connection.Connect())

	// Authenticate and create session
	authReqReply, err := clientA.connection.Request(
		context.Background(),
		"login",
		wwr.NewPayload(wwr.EncodingBinary, []byte("bla")),
	)
	require.NoError(t, err)

	session := clientA.connection.Session()
	require.Equal(t,
		session.Key, string(authReqReply.Data()),
		"Unexpected session key",
	)

	// Check status, expect 1 session, 1 connection
	require.Equal(t,
		1, server.ActiveSessionsNum(),
		"Unexpected active sessions number",
	)

	require.Equal(t,
		1, server.SessionConnectionsNum(session.Key),
		"Unexpected session connections number",
	)

	require.Len(t,
		server.SessionConnections(session.Key), 1,
		"Unexpected session connections",
	)

	// Initialize client B
	clientB := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
			// Disable autoconnect to prevent automatic reconnection
			Autoconnect: wwr.Disabled,
		},
		callbackPoweredClientHooks{},
		nil, // No TLS configuration
	)

	if err := clientB.connection.Connect(); err != nil {
		t.Fatal(err)
	}

	require.NoError(t, clientB.connection.RestoreSession(authReqReply.Data()))

	// Check status, expect 1 session, 2 connections
	require.Equal(t,
		1, server.ActiveSessionsNum(),
		"Unexpected active sessions number",
	)

	require.Equal(t,
		2, server.SessionConnectionsNum(session.Key),
		"Unexpected session connections number",
	)

	require.Len(t,
		server.SessionConnections(session.Key), 2,
		"Unexpected session connections",
	)

	// Close first connection
	_, err = clientA.connection.Request(
		context.Background(),
		"closeA",
		wwr.NewPayload(wwr.EncodingBinary, []byte("a")),
	)
	require.NoError(t, err)

	// Wait for socket to have been closed on both sides to avoid
	// a timeout on the next request
	time.Sleep(10 * time.Millisecond)

	// Test connectivity
	_, err = clientA.connection.Request(
		context.Background(),
		"testA",
		wwr.NewPayload(wwr.EncodingBinary, []byte("testA")),
	)
	require.IsType(t, wwr.DisconnectedErr{}, err)

	// Check status, expect 1 session, 1 connection
	require.Equal(t,
		1, server.ActiveSessionsNum(),
		"Unexpected active sessions number",
	)

	require.Equal(t,
		1, server.SessionConnectionsNum(session.Key),
		"Unexpected session connections number",
	)

	require.Len(t,
		server.SessionConnections(session.Key), 1,
		"Unexpected session connections",
	)

	// Close second connection
	_, err = clientB.connection.Request(
		context.Background(),
		"closeB",
		wwr.NewPayload(wwr.EncodingBinary, []byte("b")),
	)
	require.NoError(t, err)

	// Wait for socket to have been closed on both sides to avoid
	// a timeout on the next request
	time.Sleep(10 * time.Millisecond)

	// Test connectivity
	_, err = clientB.connection.Request(
		context.Background(),
		"testB",
		wwr.NewPayload(wwr.EncodingBinary, []byte("testB")),
	)
	require.IsType(t, wwr.DisconnectedErr{}, err)

	// Check status, expect 0 sessions
	require.Equal(t,
		0, server.ActiveSessionsNum(),
		"Unexpected active sessions number",
	)

	require.Equal(t,
		-1, server.SessionConnectionsNum(session.Key),
		"Unexpected session connections number",
	)

	require.Nil(t,
		server.SessionConnections(session.Key),
		"Unexpected session connections",
	)
}
