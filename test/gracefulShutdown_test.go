package test

import (
	"context"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestGracefulShutdown tests the ability of the server to delay shutdown
// until all requests and signals are processed
// and reject incoming connections and requests
// while ignoring incoming signals
//
// SIGNAL:       --->||||||||||----------------- (must finish)
// REQUEST:      ---->||||||||||---------------- (must finish)
// SRV SHUTDWN:  -------->||||||---------------- (must await req and sig)
// LATE CONN:    ---------->|------------------- (must be rejected)
// LATE REQ:     ----------->|------------------ (must be rejected)
func TestGracefulShutdown(t *testing.T) {
	expectedReqReply := []byte("ifinished")
	timeDelta := time.Duration(1)
	processesFinished := tmdwg.NewTimedWaitGroup(2, 1*time.Second)
	serverShutDown := tmdwg.NewTimedWaitGroup(
		1,
		timeDelta*500*time.Millisecond+1*time.Second,
	)

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onSignal: func(
				_ context.Context,
				_ *wwr.Client,
				_ *wwr.Message,
			) {
				time.Sleep(timeDelta * 100 * time.Millisecond)
				processesFinished.Progress(1)
			},
			onRequest: func(
				_ context.Context,
				_ *wwr.Client,
				_ *wwr.Message,
			) (wwr.Payload, error) {
				time.Sleep(timeDelta * 100 * time.Millisecond)
				return wwr.Payload{Data: expectedReqReply}, nil
			},
		},
		wwr.ServerOptions{},
	)

	serverAddr := server.Addr().String()

	// Initialize different clients for the signal,
	// the request and the late request and conn
	// to avoid serializing them because every client
	// is handled in a separate goroutine
	cltOpts := wwrclt.Options{
		DefaultRequestTimeout: 5 * time.Second,
	}
	clientSig := newCallbackPoweredClient(
		serverAddr,
		cltOpts,
		callbackPoweredClientHooks{},
	)
	clientReq := newCallbackPoweredClient(
		serverAddr,
		cltOpts,
		callbackPoweredClientHooks{},
	)
	clientLateReq := newCallbackPoweredClient(
		serverAddr,
		cltOpts,
		callbackPoweredClientHooks{},
	)

	// Disable autoconnect for the late client to enable immediate errors
	clientLateConn := newCallbackPoweredClient(
		serverAddr,
		wwrclt.Options{
			Autoconnect: wwr.Disabled,
		},
		callbackPoweredClientHooks{},
	)

	if err := clientSig.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect signal client: %s", err)
	}
	if err := clientReq.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect request client: %s", err)
	}
	if err := clientLateReq.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect late-request client: %s", err)
	}

	// Send signal and request in another parallel goroutine
	// to avoid blocking the main test goroutine when awaiting the request reply
	go func() {
		// (SIGNAL)
		if err := clientSig.connection.Signal(
			"",
			wwr.Payload{Data: []byte("test")},
		); err != nil {
			t.Errorf("Signal failed: %s", err)
		}

		// (REQUEST)
		if rep, err := clientReq.connection.Request(
			"",
			wwr.Payload{Data: []byte("test")},
		); err != nil {
			t.Errorf("Request failed: %s", err)
		} else if string(rep.Data) != string(expectedReqReply) {
			t.Errorf(
				"Expected and actual replies differ: %s | %s",
				string(expectedReqReply),
				string(rep.Data),
			)
		}
	}()

	// Request server shutdown in another parallel goroutine
	// to avoid blocking the main test goroutine when waiting
	// for the server to shut down
	go func() {
		// Wait for the signal and request to arrive and get handled,
		// then request the shutdown
		time.Sleep(timeDelta * 10 * time.Millisecond)
		// (SRV SHUTDWN)
		server.Shutdown()
		serverShutDown.Progress(1)
	}()

	// Fire late requests and late connection in another parallel goroutine
	// to avoid blocking the main test goroutine when performing them
	go func() {
		// Wait for the server to start shutting down
		time.Sleep(timeDelta * 20 * time.Millisecond)

		// Verify connection establishment during shutdown (LATE CONN)
		if err := clientLateConn.connection.Connect(); err == nil {
			t.Errorf("Expected late connection to be rejected, " +
				"though it still was accepted",
			)
		}

		// Verify request rejection during shutdown (LATE REQ)
		_, lateReqErr := clientLateReq.connection.Request(
			"",
			wwr.Payload{Data: []byte("test")},
		)
		switch err := lateReqErr.(type) {
		case wwr.ReqSrvShutdownErr:
			break
		case wwr.ReqErr:
			t.Errorf("Expected special server shutdown error, "+
				"got regular request error: %s",
				err,
			)
		default:
			t.Errorf("Expected request during shutdown to be rejected " +
				"with special error type",
			)
		}

		processesFinished.Progress(1)
	}()

	// Await server shutdown, timeout if necessary
	if err := serverShutDown.Wait(); err != nil {
		t.Fatalf("Expected server to shut down within n seconds")
	}

	// Expect both the signal and the request to have completed properly
	if err := processesFinished.Wait(); err != nil {
		t.Fatalf("Expected signal and request to have finished processing")
	}
}
