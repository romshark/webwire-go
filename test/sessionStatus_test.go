package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestSessionStatus tests session monitoring methods
func TestSessionStatus(t *testing.T) {
	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				clt *wwr.Client,
				_ *wwr.Message,
			) (wwr.Payload, error) {
				// Try to create a new session
				if err := clt.CreateSession(nil); err != nil {
					return wwr.Payload{}, err
				}

				// Return the key of the newly created session
				// (use default binary encoding)
				return wwr.Payload{
					Data: []byte(clt.SessionKey()),
				}, nil
			},
		},
		wwr.ServerOptions{},
	)

	actSess := server.ActiveSessionsNum()
	if actSess != 0 {
		t.Fatalf("Unexpected number of active sessions: %d", actSess)
	}

	// Initialize client A
	clientA := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	// Authenticate and create session
	authReqReply, err := clientA.connection.Request("login", wwr.Payload{
		Data: []byte("bla"),
	})
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	session := clientA.connection.Session()
	if session.Key != string(authReqReply.Data) {
		t.Fatalf("Unexpected session key")
	}

	// Check status, expect 1 session with 1 connection
	actSess = server.ActiveSessionsNum()
	if actSess != 1 {
		t.Fatalf("Unexpected active sessions number: %d", actSess)
	}

	sessConnsNum := server.SessionConnectionsNum(session.Key)
	if sessConnsNum != 1 {
		t.Fatalf("Unexpected session connections number: %d", sessConnsNum)
	}

	sessConns := server.SessionConnections(session.Key)
	if len(sessConns) != 1 {
		t.Fatalf("Unexpected session connections: %d", len(sessConns))
	}

	// Initialize client B
	clientB := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	if err := clientB.connection.RestoreSession(authReqReply.Data); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	// Check status, expect 1 session with 2 connections
	actSess = server.ActiveSessionsNum()
	if actSess != 1 {
		t.Fatalf("Unexpected active sessions number: %d", actSess)
	}

	sessConnsNum = server.SessionConnectionsNum(session.Key)
	if sessConnsNum != 2 {
		t.Fatalf("Unexpected session connections number: %d", sessConnsNum)
	}

	sessConns = server.SessionConnections(session.Key)
	if len(sessConns) != 2 {
		t.Fatalf("Unexpected session connections: %d", len(sessConns))
	}

	// Close first connection
	if err := clientA.connection.CloseSession(); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	// Check status, expect 1 session with 1 connection
	actSess = server.ActiveSessionsNum()
	if actSess != 1 {
		t.Fatalf("Unexpected active sessions number: %d", actSess)
	}

	sessConnsNum = server.SessionConnectionsNum(session.Key)
	if sessConnsNum != 1 {
		t.Fatalf("Unexpected session connections number: %d", sessConnsNum)
	}

	sessConns = server.SessionConnections(session.Key)
	if len(sessConns) != 1 {
		t.Fatalf("Unexpected session connections: %d", len(sessConns))
	}

	// Close second connection
	if err := clientB.connection.CloseSession(); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	// Check status, expect 0 sessions
	actSess = server.ActiveSessionsNum()
	if actSess != 0 {
		t.Fatalf("Unexpected active sessions number: %d", actSess)
	}

	sessConnsNum = server.SessionConnectionsNum(session.Key)
	if sessConnsNum != -1 {
		t.Fatalf("Unexpected session connections number: %d", sessConnsNum)
	}

	sessConns = server.SessionConnections(session.Key)
	if sessConns != nil {
		t.Fatalf("Unexpected session connections: %d", len(sessConns))
	}
}
