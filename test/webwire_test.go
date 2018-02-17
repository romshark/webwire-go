package test

import (
	"testing"
	"os"
	"fmt"
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
		onSaveSession, onFindSession, onSessionClosure,
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
		t.Errorf("Invalid %s: payload doesn't match:\n expected: '%s'\n actual:   '%s'",
			name,
			string(expected),
			string(actual),
		)
	}
}

func compareSessions(t *testing.T, expected, actual *webwire.Session) {
	if actual == nil && expected == nil {
		return
	} else if (actual == nil && expected != nil) || (expected == nil && actual != nil) {
		t.Errorf("Sessions differ:\n expected: '%v'\n actual:   '%v'",
			expected,
			actual,
		)
	} else if actual.Key != expected.Key {
		t.Errorf("Session keys differ:\n expected: '%s'\n actual:   '%s'",
			expected.Key,
			actual.Key,
		)
	}
	// TODO: get session creation time from actual server time
	/*
	else if actual.CreationDate != expected.CreationDate {
		t.Errorf("Session creation dates differ:\n expected: '%s'\n actual:   '%s'",
			expected.CreationDate,
			actual.CreationDate,
		)
	}
	*/
	// TODO: implement user agent, OS and info comparison
	/*
	else if actual.UserAgent != expected.UserAgent {
		t.Errorf("Session user agents differ:\n expected: '%s'\n actual:   '%s'",
			expected.UserAgent,
			actual.UserAgent,
		)
	} else if actual.OperatingSystem != expected.OperatingSystem {
		t.Errorf("Session operating systems differ:\n expected: '%v'\n actual:   '%v'",
			expected.OperatingSystem,
			actual.OperatingSystem,
		)
	} else if !reflect.DeepEqual(expected.Info, actual.Info) {
		t.Errorf("Session info differs:\n expected: '%v'\n actual:   '%v'",
			expected.Info,
			actual.Info,
		)
	}
	*/
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
			comparePayload(t, "client request", expectedRequestPayload, msg.Payload)
			return expectedReplyPayload, nil
		},
		nil, nil, nil, nil,
	)
	go server.Run()

	// Initialize client
	client := webwire_client.NewClient(
		server.Addr,
		nil,
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
			comparePayload(t, "client signal", expectedSignalPayload, msg.Payload)

			// Synchronize, notify signal arival
			wait <- true
		},
		nil, nil, nil, nil, nil,
	)
	go server.Run()

	// Initialize client
	client := webwire_client.NewClient(
		server.Addr,
		nil,
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

			// Send signal
			if err := client.Signal(expectedSignalPayload); err != nil {
				t.Fatalf("Couldn't send signal to client: %s", err)
			}
		}, nil, nil, nil, nil, nil, nil)
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
		nil,
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

// TestSessionCreation verifies the server is connectable,
// and is able to receives requests and create sessions for the requesting client
func TestSessionCreation(t *testing.T) {
	var finish sync.WaitGroup
	var createdSession *webwire.Session
	finish.Add(2)

	// Initialize webwire server
	server := setupServer(
		t,
		nil,
		nil,
		// onRequest
		func(ctx context.Context) ([]byte, *webwire.Error) {
			defer finish.Done()

			// Extract request message and requesting client from the context
			msg := ctx.Value(webwire.MESSAGE).(webwire.Message)

			// Create a new session
			newSession := webwire.NewSession(
				webwire.Os_UNKNOWN,
				"user agent",
				nil,
			)
			createdSession = &newSession

			// Try to register the newly created session and bind it to the client
			if err := msg.Client.CreateSession(createdSession); err != nil {
				return nil, &webwire.Error {
					"INTERNAL_ERROR",
					fmt.Sprintf("Internal server error: %s", err),
				}
			}

			// Return the key of the newly created session
			return []byte(createdSession.Key), nil
		},
		// OnSaveSession
		func(session *webwire.Session) error {
			// Verify the session
			compareSessions(t, createdSession, session)
			return nil
		},
		// OnFindSession
		func(_ string) (*webwire.Session, error) {
			return nil, nil
		},
		// OnSessionClosure
		func(_ string) error {
			return nil
		},
		nil,
	)
	go server.Run()

	// Initialize client
	client := webwire_client.NewClient(
		server.Addr,
		nil,
		// On session creation
		func (newSession *webwire.Session) {
			defer finish.Done()

			// Verify reply
			compareSessions(t, createdSession, newSession)
		},
		5 * time.Second,
		os.Stdout,
		os.Stderr,
	)
	defer client.Close()

	// Send request and await reply
	reply, err := client.Request([]byte("credentials"))
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify reply
	comparePayload(t, "server reply", []byte(createdSession.Key), reply)

	// Verify client session
	finish.Wait()
}


// TestAuthentication verifies the server is connectable,
// and is able to receives requests and signals, create sessions
// and identify clients during request- and signal handling
func TestAuthentication(t *testing.T) {
	var clientSignal sync.WaitGroup
	clientSignal.Add(1)
	var createdSession *webwire.Session
	expectedCredentials := []byte("secret_credentials")
	expectedConfirmation := []byte("session_is_correct")
	currentStep := 1

	// Initialize webwire server
	server := setupServer(
		t,
		nil,
		// onSignal
		func(ctx context.Context) {
			defer clientSignal.Done()
			// Extract request message and requesting client from the context
			msg := ctx.Value(webwire.MESSAGE).(webwire.Message)
			compareSessions(t, createdSession, msg.Client.Session)
		},
		// onRequest
		func(ctx context.Context) ([]byte, *webwire.Error) {
			// Extract request message and requesting client from the context
			msg := ctx.Value(webwire.MESSAGE).(webwire.Message)

			// If already authenticated then check session
			if currentStep > 1 {
				compareSessions(t, createdSession, msg.Client.Session)
				return expectedConfirmation, nil
			}

			// Create a new session
			newSession := webwire.NewSession(
				webwire.Os_UNKNOWN,
				"user agent",
				nil,
			)
			createdSession = &newSession

			// Try to register the newly created session and bind it to the client
			if err := msg.Client.CreateSession(createdSession); err != nil {
				return nil, &webwire.Error {
					"INTERNAL_ERROR",
					fmt.Sprintf("Internal server error: %s", err),
				}
			}

			// Authentication step is passed
			currentStep = 2

			// Return the key of the newly created session
			return []byte(createdSession.Key), nil
		},
		// OnSaveSession
		func(session *webwire.Session) error {
			// Verify the session
			compareSessions(t, createdSession, session)
			return nil
		},
		// OnFindSession
		func(_ string) (*webwire.Session, error) {
			return nil, nil
		},
		// OnSessionClosure
		func(_ string) error {
			return nil
		},
		nil,
	)
	go server.Run()

	// Initialize client
	client := webwire_client.NewClient(
		server.Addr,
		nil,
		nil,
		5 * time.Second,
		os.Stdout,
		os.Stderr,
	)
	defer client.Close()

	// Send authentication request and await reply
	authReqReply, err := client.Request(expectedCredentials)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify reply
	comparePayload(t, "authentication reply", []byte(createdSession.Key), authReqReply)


	// Send a test-request to verify the session on the server and await response
	testReqReply, err := client.Request(expectedCredentials)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify reply
	comparePayload(t, "test reply", expectedConfirmation, testReqReply)

	// Send a test-signal to verify the session on the server
	if err := client.Signal(expectedCredentials); err != nil {
		t.Fatalf("Request failed: %s", err)
	}
	clientSignal.Wait()
}
