package test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/assert"
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
	releaseReqHandler := sync.WaitGroup{}
	releaseReqHandler.Add(1)
	releaseSigHandler := sync.WaitGroup{}
	releaseSigHandler.Add(1)
	firstReqAndSigReceived := sync.WaitGroup{}
	firstReqAndSigReceived.Add(2)
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
				defer handlersFinished.Done()

				if string(msg.Name()) == "1" {
					firstReqAndSigReceived.Done()
				}
				// Sleep after the first signal was marked as done
				releaseSigHandler.Wait()
			},
			Request: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				defer handlersFinished.Done()

				if string(msg.Name()) == "1" {
					firstReqAndSigReceived.Done()
				}
				releaseReqHandler.Wait()
				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize different clients for the signal,
	// the request and the late request and conn
	// to avoid serializing them because every client
	// is handled in a separate goroutine
	clientSig, _ := setup.NewClientSocket()
	clientReq, _ := setup.NewClientSocket()
	clientLateReq, _ := setup.NewClientSocket()

	// Disable autoconnect for the late client to enable immediate errors
	clientLateConn, err := setup.NewDisconnectedClientSocket()
	require.NoError(t, err)

	// Send signal and request in another parallel goroutine
	// to avoid blocking the main test goroutine when awaiting the request reply
	go func() {
		// (SIGNAL)
		signal(t, clientSig, []byte("1"), payload.Payload{})

		// (REQUEST)
		requestSuccess(t, clientReq, 32, []byte("1"), payload.Payload{})
	}()

	// Request server shutdown in another parallel goroutine
	// to avoid blocking the main test goroutine when waiting
	// for the server to shut down
	go func() {
		// Wait for the signal and request to arrive and get handled,
		// then request the shutdown
		firstReqAndSigReceived.Wait()

		// Wait a little before claming that the server is shutting down
		time.AfterFunc(50*time.Millisecond, func() {
			serverShuttingDown.Done()
		})

		// (SRV SHUTDOWN)
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
			clientLateConn.Dial(time.Time{}),
			"Expected late connection to be rejected, "+
				"though it still was accepted",
		)

		// Verify request rejection during shutdown (LATE REQ)
		rep := request(t, clientLateReq, 32, []byte("r"), payload.Payload{})
		assert.Equal(t, message.MsgReplyShutdown, rep.MsgType)
		assert.Nil(t, rep.MsgPayload.Data)

		// Release the handlers to allow the server to finally shutdown
		releaseSigHandler.Done()
		releaseReqHandler.Done()
	}()

	// Await actual server shutdown
	serverShutDown.Wait()

	// Expect both the signal and the request handlers to have properly finished
	handlersFinished.Wait()
}
