package test

import (
	"testing"
	"os"
	"net"
	"time"
	"sync"
	"context"
	"reflect"

	webwire "github.com/qbeon/webwire-go"
	webwire_client "github.com/qbeon/webwire-go/client"
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

// setupServer helps setting up and launching the server together with the hosting http server
func setupServer(
	t *testing.T,
	onClientConnected webwire.OnClientConnected,
	onSignal webwire.OnSignal,
	onRequest webwire.OnRequest,
	onSessionCreation webwire.OnSessionCreation,
	onSaveSession webwire.OnSaveSession,
	onFindSession webwire.OnFindSession,
	onSessionClosure webwire.OnSessionClosure,
	onCORS webwire.OnCORS,
) (srv *webwire.Server) {
	// Initialize webwire server
	webwireServer, err := webwire.NewServer(
		"127.0.0.1:0",
		onClientConnected,
		onSignal, onRequest,
		onSessionCreation, onSaveSession, onFindSession, onSessionClosure,
		onCORS,
		os.Stdout, os.Stderr,
	)
	if err != nil {
		t.Fatalf("Failed creating a new WebWire server instance: %s", err)
	}

	return webwireServer
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
	server := setupServer(
		t,
		nil, nil,
		func(ctx context.Context) ([]byte, *webwire.Error) {
			// Extract request message from the context
			msg := ctx.Value(webwire.MESSAGE).(webwire.Message)

			// Verify request payload
			comparePayload(t, "client request", expectedRequestPayload, msg.Payload())
			return expectedReplyPayload, nil
		},
		nil, nil, nil, nil, nil,
	)
	go server.Run()

	// Initialize client
	client := webwire_client.NewClient(
		server.Addr,
		nil,
		5 * time.Second,
		os.Stdout,
		os.Stderr,
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
	server := setupServer(
		t,
		nil,
		func(ctx context.Context) {
			// Extract signal message from the context
			msg := ctx.Value(webwire.MESSAGE).(webwire.Message)

			// Verify signal payload
			comparePayload(t, "client signal", expectedSignalPayload, msg.Payload())

			// Synchronize, notify signal arival
			wait <- true
		},
		nil, nil, nil, nil, nil, nil,
	)
	go server.Run()

	// Initialize client
	client := webwire_client.NewClient(
		server.Addr,
		nil,
		5 * time.Second,
		os.Stdout,
		os.Stderr,
	)

	// Send signal
	err := client.Signal(expectedSignalPayload)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Synchronize, await signal arival
	<- wait
}

// TestServerSignal verifies the server is connectable and sends signals correctly
func TestServerSignal(t *testing.T) {
	expectedSignalPayload := []byte("webwire_test_SERVER_SIGNAL_payload")
	var addr string
	var server *webwire.Server
	var finish sync.WaitGroup
	finish.Add(1)
	initClient := make(chan bool, 1)
	sendSignal := make(chan bool, 1)

	// Initialize webwire server
	go func() {
		server = setupServer(t, func (client *webwire.Client) {

			// Verify client is listed
			/*
			if server.ClientsNum() != 1 {
				finish.Done()
				t.Fatalf(
					"Unexpected list of connected clients (%d), expected 1 client to be connected",
					server.ClientsNum(),
				)
			}
			*/

			// Send signal*
			if err := client.Signal(expectedSignalPayload); err != nil {
				t.Fatalf("Couldn't send signal to client: %s", err)
			}
		}, nil, nil, nil, nil, nil, nil, nil)
		go server.Run()
		addr = server.Addr

		// Synchronize, initialize client
		initClient <- true

		// Synchronize, wait for the client to launch and require the signal to be sent
		<- sendSignal
	}()

	// Synchronize, await server initialization
	<- initClient

	// Initialize client
	client := webwire_client.NewClient(
		addr,
		func(signalPayload []byte) {
			// Verify server signal payload
			comparePayload(t, "server signal", expectedSignalPayload, signalPayload)

			// Synchronize, unlock main goroutine to pass the test case
			finish.Done()
		},
		5 * time.Second,
		os.Stdout,
		os.Stderr,
	)
	defer client.Close()

	// Connect client
	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect client: %s", err)
	}

	// Synchronize, notify the server the client was initialized and request the signal
	sendSignal <- true

	// Synchronize, await signal arival
	finish.Wait()
}
