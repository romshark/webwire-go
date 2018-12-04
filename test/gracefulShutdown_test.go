package test

import (
	"context"
	"sync"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGracefulShutdown tests the ability of the server to delay shutdown
// until all requests and signals are processed
// and reject incoming connections and requests
// while ignoring incoming signals
//
// SIGNAL:       |-->||||||||||----------------- (must finish)
// REQUEST:      ||-->||||||||||---------------- (must finish)
// SRV SHUTDWN:  |||----->||||||---------------- (must await req and sig)
// LATE CONN:    |||------->|------------------- (must be rejected)
// LATE REQ:     ||||------->|------------------ (must be rejected)
func TestGracefulShutdown(t *testing.T) {
	expectedReqReply := []byte("i_finished")
	handlerExecutionDuration := 100 * time.Millisecond
	firstReqAndSigSent := sync.WaitGroup{}
	firstReqAndSigSent.Add(2)
	serverShuttingDown := sync.WaitGroup{}
	serverShuttingDown.Add(1)
	handlersFinished := sync.WaitGroup{}
	handlersFinished.Add(2)
	serverShutDown := sync.WaitGroup{}
	serverShutDown.Add(1)

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Signal: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) {
				if string(msg.Name()) == "1" {
					firstReqAndSigSent.Done()
				}
				// Sleep after the first signal was marked as done
				time.Sleep(handlerExecutionDuration)
				handlersFinished.Done()
			},
			Request: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				if string(msg.Name()) == "1" {
					firstReqAndSigSent.Done()
				}
				time.Sleep(handlerExecutionDuration)
				return wwr.Payload{Data: expectedReqReply}, nil
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize different clients for the signal,
	// the request and the late request and conn
	// to avoid serializing them because every client
	// is handled in a separate goroutine
	cltOpts := wwrclt.Options{
		DefaultRequestTimeout: 5 * time.Second,
		Autoconnect:           wwr.Disabled,
	}
	clientSig := setup.NewClient(
		cltOpts,
		nil, // Use the default transport implementation
		TestClientHooks{},
	)
	clientReq := setup.NewClient(
		cltOpts,
		nil, // Use the default transport implementation
		TestClientHooks{},
	)
	clientLateReq := setup.NewClient(
		cltOpts,
		nil, // Use the default transport implementation
		TestClientHooks{},
	)

	require.NoError(t, clientSig.Connection.Connect())
	require.NoError(t, clientReq.Connection.Connect())
	require.NoError(t, clientLateReq.Connection.Connect())

	// Disable autoconnect for the late client to enable immediate errors
	clientLateConn := setup.NewClient(
		wwrclt.Options{
			Autoconnect: wwr.Disabled,
		},
		nil, // Use the default transport implementation
		TestClientHooks{},
	)

	// Send signal and request in another parallel goroutine
	// to avoid blocking the main test goroutine when awaiting the request reply
	go func() {
		// (SIGNAL)
		assert.NoError(t, clientSig.Connection.Signal(
			context.Background(),
			[]byte("1"),
			wwr.Payload{Data: []byte("test")},
		))

		// (REQUEST)
		rep, err := clientReq.Connection.Request(
			context.Background(),
			[]byte("1"),
			wwr.Payload{Data: []byte("test")},
		)
		assert.NoError(t, err)
		assert.Equal(t, string(rep.Payload()), string(expectedReqReply))
		handlersFinished.Done()
		rep.Close()
	}()

	// Request server shutdown in another parallel goroutine
	// to avoid blocking the main test goroutine when waiting
	// for the server to shut down
	go func() {
		// Wait for the signal and request to arrive and get handled,
		// then request the shutdown
		firstReqAndSigSent.Wait()

		// (SRV SHUTDWN)
		serverShuttingDown.Done()
		setup.Server.Shutdown()
		serverShutDown.Done()
	}()

	// Wait for the server to start shutting down and fire late requests
	// and late connection in another parallel goroutine
	// to avoid blocking the main test goroutine when performing them
	go func() {
		// Wait for the server to start shutting down
		serverShuttingDown.Wait()

		// Verify connection establishment during shutdown (LATE CONN)
		assert.Error(t,
			clientLateConn.Connection.Connect(),
			"Expected late connection to be rejected, "+
				"though it still was accepted",
		)

		// Verify request rejection during shutdown (LATE REQ)
		_, lateReqErr := clientLateReq.Connection.Request(
			context.Background(),
			nil,
			wwr.Payload{Data: []byte("test")},
		)
		switch err := lateReqErr.(type) {
		case wwr.ServerShutdownErr:
			break
		case wwr.RequestErr:
			t.Errorf("Expected special server shutdown error, "+
				"got regular request error: %s",
				err,
			)
		default:
			t.Errorf("Expected request during shutdown to be rejected " +
				"with special error type",
			)
		}
	}()

	// Await server shutdown, timeout if necessary
	serverShutDown.Wait()

	// Expect both the signal and the request to have completed properly
	handlersFinished.Wait()
}
