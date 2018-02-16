package webwire

import (
	"testing"
	"os"
	"fmt"
	"net"
	"time"
	"context"
	"reflect"
	"net/http"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// createHttpServer helps setting up the HTTP server basement for the webwire server
func createHttpServer(handler http.Handler) (srv *http.Server, addr string, err error) {
	httpServer := &http.Server{Addr: "127.0.0.1:0", Handler: handler}
	addr = httpServer.Addr
	if addr == "" {
		addr = ":http"
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, "", err
	}

	// Launch server
	go func() {
		err = httpServer.Serve(tcpKeepAliveListener{listener.(*net.TCPListener)})
		if err != nil {
			panic(fmt.Errorf("HTTP Server failure: %s", err))
		}
	}()

	return httpServer, listener.Addr().String(), nil
}

// setup helps setting up and launching everything, the client, the server and the http server.
// The returned teardown handle must be defered right after calling setup for proper shutdown
func setup(
	t *testing.T,
	onSignal OnSignal,
	onRequest OnRequest,
	onSessionCreation OnSessionCreation,
	onSaveSession OnSaveSession,
	onFindSession OnFindSession,
	onSessionClosure OnSessionClosure,
	onCORS OnCORS,
) (srv *Server, httpSrv *http.Server, clt *Client) {
	// Initialize webwire server
	webwireServer := NewServer(
		onSignal, onRequest,
		onSessionCreation, onSaveSession, onFindSession, onSessionClosure,
		onCORS,
		os.Stdout, os.Stderr,
	)

	httpSrv, addr, err := createHttpServer(webwireServer)
	if err != nil {
		t.Fatalf("Failed setting up HTTP server: %s", err)
	}

	// Initialize client
	client := NewClient(addr, 5 * time.Second, os.Stdout, os.Stderr)

	return &webwireServer, httpSrv, &client
}

func comparePayload(t *testing.T, name string, expected, actual []byte) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Invalid %s: payload doesn't match:\n expected: '%s'\n actual:   '%s'",
			name,
			string(expected),
			string(actual),
		)
	}
}

// TESTS

// TestClientRequest verifies the server is connectable,
// receives requests and answers them correctly
func TestClientRequest(t *testing.T) {
	expectedRequestPayload := []byte("webwire_test_REQUEST_payload")
	expectedReplyPayload := []byte("webwire_test_RESPONSE_message")

	// Initialize webwire server given only the request
	_, _, client := setup(
		t,
		nil,
		func(ctx context.Context) ([]byte, *Error) {
			// Extract request message from the context
			msg := ctx.Value(MESSAGE).(Message)

			// Verify request payload
			comparePayload(t, "client request", expectedRequestPayload, msg.Payload())
			return expectedReplyPayload, nil
		},
		nil, nil, nil, nil, nil,
	)

	// Send request and await reply
	reply, err := client.Request(expectedRequestPayload)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify reply
	comparePayload(t, "server reply", expectedReplyPayload, reply)
}

// TestClientSignal verifies the server is connectable and receives signals correctly
func TestClientSignal(t *testing.T) {
	expectedSignalPayload := []byte("webwire_test_SIGNAL_payload")
	wait := make(chan bool)

	// Initialize webwire server given only the signal handler
	_, _, client := setup(
		t,
		func(ctx context.Context) {
			// Extract signal message from the context
			msg := ctx.Value(MESSAGE).(Message)

			// Verify signal payload
			comparePayload(t, "client signal", expectedSignalPayload, msg.Payload())

			// Notify signal arival
			wait <- true
		},
		nil, nil, nil, nil, nil, nil,
	)

	// Send signal
	err := client.Signal(expectedSignalPayload)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Await signal arival
	<- wait
}
