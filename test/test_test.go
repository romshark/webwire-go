package test

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/qbeon/webwire-go/transport/memchan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type serverSetup struct {
	Transport wwr.Transport
	Server    wwr.Server
}

type testServerSetup struct {
	t *testing.T
	serverSetup
}

// setupServer helps setting up and launching the server together with the
// underlying transport
func setupServer(
	impl *serverImpl,
	opts wwr.ServerOptions,
	trans wwr.Transport,
) (serverSetup, error) {
	// Setup headed server on arbitrary port
	if impl.onClientConnected == nil {
		impl.onClientConnected = func(_ wwr.Connection) {}
	}
	if impl.onClientDisconnected == nil {
		impl.onClientDisconnected = func(_ wwr.Connection, _ error) {}
	}
	if impl.onSignal == nil {
		impl.onSignal = func(
			_ context.Context,
			_ wwr.Connection,
			_ wwr.Message,
		) {
		}
	}
	if impl.onRequest == nil {
		impl.onRequest = func(
			_ context.Context,
			_ wwr.Connection,
			_ wwr.Message,
		) (response wwr.Payload, err error) {
			return wwr.Payload{}, nil
		}
	}

	// Use default session manager if no specific one is defined
	if opts.SessionManager == nil {
		opts.SessionManager = newInMemSessManager()
	}

	// Use default address
	opts.Host = "127.0.0.1:0"

	// Use the transport layer implementation specified by the CLI arguments
	if trans == nil {
		// Use default configuration
		trans = &memchan.Transport{}
	}

	// Initialize webwire server
	server, err := wwr.NewServer(impl, opts, trans)
	if err != nil {
		return serverSetup{}, fmt.Errorf(
			"failed setting up server instance: %s",
			err,
		)
	}

	// Run server in a separate goroutine
	go func() {
		if err := server.Run(); err != nil {
			panic(fmt.Errorf("server failed: %s", err))
		}
	}()

	// Return reference to the server and the address its bound to
	return serverSetup{
		Server:    server,
		Transport: trans,
	}, nil
}

// setupTestServer creates a new server setup failing the test immediately if
// the anything went wrong
func setupTestServer(
	t *testing.T,
	impl *serverImpl,
	opts wwr.ServerOptions,
	trans wwr.Transport,
) testServerSetup {
	setup, err := setupServer(impl, opts, trans)
	require.NoError(t, err)
	return testServerSetup{t, setup}
}

// newClient sets up a new test client instance
func (setup *serverSetup) newClient(
	options wwrclt.Options,
	transport wwr.ClientTransport,
	hooks testClientHooks,
) (*testClient, error) {
	return newClient(
		setup.Server.Address(),
		options,
		transport,
		hooks,
	)
}

// newClient sets up a new test client instance
func (setup *testServerSetup) newClient(
	options wwrclt.Options,
	transport wwr.ClientTransport,
	hooks testClientHooks,
) *testClient {
	if transport == nil {
		transport = &memchan.ClientTransport{
			Server: setup.Transport.(*memchan.Transport),
		}
	}

	clt, err := newClient(
		setup.Server.Address(),
		options,
		transport,
		hooks,
	)
	require.NoError(setup.t, err)
	return clt
}

// newClient sets up a new test client instance
func newClient(
	serverAddr url.URL,
	options wwrclt.Options,
	clientTransport wwr.ClientTransport,
	hooks testClientHooks,
) (*testClient, error) {
	newClt := &testClient{
		hooks: hooks,
	}

	// Initialize connection
	conn, err := wwrclt.NewClient(serverAddr, newClt, options, clientTransport)
	if err != nil {
		return nil, fmt.Errorf("failed setting up client instance: %s", err)
	}

	newClt.connection = conn

	return newClt, nil
}

// newClientSocket creates a new raw client socket connected to the server
func (setup *serverSetup) newClientSocket() (wwr.Socket, error) {
	switch srvTrans := setup.Transport.(type) {
	case *memchan.Transport:
		_, sock := memchan.NewEntangledSockets(srvTrans)
		if err := sock.Dial(url.URL{}, time.Time{}); err != nil {
			return nil, fmt.Errorf("memchan dial failed: %s", err)
		}
		return sock, nil
	}
	return nil, fmt.Errorf(
		"unexpected server transport implementation: %s",
		reflect.TypeOf(setup.Transport),
	)
}

// newClientSocket creates a new raw client socket connected to the server
func (setup *testServerSetup) newClientSocket() wwr.Socket {
	sock, err := setup.serverSetup.newClientSocket()
	require.NoError(setup.t, err)
	return sock
}

func compareSessions(t *testing.T, expected, actual *wwr.Session) {
	if actual == nil && expected == nil {
		return
	}

	assert.NotNil(t, expected)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Key, actual.Key)
	assert.Equal(t, expected.Creation.Unix(), actual.Creation.Unix())
}
