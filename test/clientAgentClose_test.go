package test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientAgentClose tests closing a client agent on the server-side
func TestClientAgentClose(t *testing.T) {
	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				clt *wwr.Client,
				msg wwr.Message,
			) (wwr.Payload, error) {
				switch msg.Name() {
				case "closeA":
					fallthrough
				case "closeB":
					clt.Close()
					return nil, nil
				case "login":
					// Try to create a new session
					if err := clt.CreateSession(nil); err != nil {
						return nil, err
					}

					// Return the key of the newly created session
					// (use default binary encoding)
					return wwr.NewPayload(
						wwr.EncodingBinary,
						[]byte(clt.SessionKey()),
					), nil
				}
				return nil, fmt.Errorf("Invalid request %s", msg.Name())
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
			// Disable autoconnect to prevent automatic reconnection
			Autoconnect: wwr.Disabled,
		},
		callbackPoweredClientHooks{},
	)

	if err := clientA.connection.Connect(); err != nil {
		t.Fatal(err)
	}

	// Authenticate and create session
	authReqReply, err := clientA.connection.Request("login", wwr.NewPayload(
		wwr.EncodingBinary,
		[]byte("bla"),
	))
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	session := clientA.connection.Session()
	if session.Key != string(authReqReply.Data()) {
		t.Fatalf("Unexpected session key")
	}

	// Check status, expect 1 session, 1 connection
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
			// Disable autoconnect to prevent automatic reconnection
			Autoconnect: wwr.Disabled,
		},
		callbackPoweredClientHooks{},
	)

	if err := clientB.connection.Connect(); err != nil {
		t.Fatal(err)
	}

	if err := clientB.connection.RestoreSession(
		authReqReply.Data(),
	); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	// Check status, expect 1 session, 2 connections
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
	if _, err := clientA.connection.Request("closeA", wwr.NewPayload(
		wwr.EncodingBinary,
		[]byte("a"),
	)); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Wait for socket to have been closed on both sides to avoid
	// a timeout on the next request
	time.Sleep(10 * time.Millisecond)

	// Test connectivity
	_, err = clientA.connection.Request(
		"testA",
		wwr.NewPayload(wwr.EncodingBinary, []byte("testA")),
	)
	if _, isDisconnErr := err.(wwr.DisconnectedErr); !isDisconnErr {
		t.Fatalf(
			"Expected disconnected error, got: %s | %s",
			reflect.TypeOf(err),
			err,
		)
	}

	// Check status, expect 1 session, 1 connection
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
	if _, err := clientB.connection.Request("closeB", wwr.NewPayload(
		wwr.EncodingBinary,
		[]byte("b"),
	)); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Wait for socket to have been closed on both sides to avoid
	// a timeout on the next request
	time.Sleep(10 * time.Millisecond)

	// Test connectivity
	_, err = clientB.connection.Request(
		"testB",
		wwr.NewPayload(wwr.EncodingBinary, []byte("testB")),
	)
	if _, isDisconnErr := err.(wwr.DisconnectedErr); !isDisconnErr {
		t.Fatalf(
			"Expected disconnected error, got: %s | %s",
			reflect.TypeOf(err),
			err,
		)
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
