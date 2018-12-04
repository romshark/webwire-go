package test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/qbeon/webwire-go/message"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/qbeon/webwire-go/transport/memchan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ServerSetup struct {
	Transport wwr.Transport
	Server    wwr.Server
}

type TestServerSetup struct {
	t *testing.T
	ServerSetup
}

// SetupServer helps setting up and launching the server together with the
// underlying transport
func SetupServer(
	impl *ServerImpl,
	opts wwr.ServerOptions,
	trans wwr.Transport,
) (ServerSetup, error) {
	// Use default session manager if no specific one is defined
	if opts.SessionManager == nil {
		opts.SessionManager = newInMemSessManager()
	}

	// Use the transport layer implementation specified by the CLI arguments
	if trans == nil {
		// Use default configuration
		trans = &memchan.Transport{}
	}

	// Initialize webwire server
	server, err := wwr.NewServer(impl, opts, trans)
	if err != nil {
		return ServerSetup{}, fmt.Errorf(
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
	return ServerSetup{
		Server:    server,
		Transport: trans,
	}, nil
}

// SetupTestServer creates a new server setup failing the test immediately if
// the anything went wrong
func SetupTestServer(
	t *testing.T,
	impl *ServerImpl,
	opts wwr.ServerOptions,
	trans wwr.Transport,
) TestServerSetup {
	setup, err := SetupServer(impl, opts, trans)
	require.NoError(t, err)
	return TestServerSetup{t, setup}
}

// NewClient sets up a new test client instance
func (setup *ServerSetup) NewClient(
	options wwrclt.Options,
	transport wwr.ClientTransport,
	hooks TestClientHooks,
) (*TestClient, error) {
	return NewClient(
		options,
		transport,
		hooks,
	)
}

// NewClient sets up a new test client instance
func (setup *TestServerSetup) NewClient(
	options wwrclt.Options,
	transport wwr.ClientTransport,
	hooks TestClientHooks,
) *TestClient {
	if transport == nil {
		transport = &memchan.ClientTransport{
			Server: setup.Transport.(*memchan.Transport),
		}
	}

	clt, err := NewClient(
		options,
		transport,
		hooks,
	)
	require.NoError(setup.t, err)
	return clt
}

// NewClient sets up a new test client instance
func NewClient(
	options wwrclt.Options,
	clientTransport wwr.ClientTransport,
	hooks TestClientHooks,
) (*TestClient, error) {
	newClt := &TestClient{Hooks: hooks}

	// Initialize connection
	conn, err := wwrclt.NewClient(newClt, options, clientTransport)
	if err != nil {
		return nil, fmt.Errorf("failed setting up client instance: %s", err)
	}

	newClt.Connection = conn

	return newClt, nil
}

// NewDisconnectedClientSocket creates a new raw disconnected client socket
func (setup *ServerSetup) NewDisconnectedClientSocket() (
	wwr.ClientSocket,
	error,
) {
	var sock wwr.ClientSocket
	switch srvTrans := setup.Transport.(type) {
	case *memchan.Transport:
		_, sock = memchan.NewEntangledSockets(srvTrans)
		return sock, nil
	}
	return nil, fmt.Errorf(
		"unexpected server transport implementation: %s",
		reflect.TypeOf(setup.Transport),
	)
}

// NewClientSocket creates a new raw client socket connected to the server
func (setup *ServerSetup) NewClientSocket() (
	wwr.Socket,
	message.ServerConfiguration,
	error,
) {
	sock, err := setup.NewDisconnectedClientSocket()
	if err != nil {
		return nil, message.ServerConfiguration{}, err
	}

	// Establish a connection
	if err := sock.Dial(time.Time{}); err != nil {
		return nil, message.ServerConfiguration{}, fmt.Errorf(
			"memchan dial failed: %s",
			err,
		)
	}

	// Read the server configuration push message
	msg := message.NewMessage(32)
	if err := sock.Read(msg, time.Time{}); err != nil {
		return nil, message.ServerConfiguration{}, fmt.Errorf(
			"couldn't read server configuration push message: %s",
			err,
		)
	}

	return sock, msg.ServerConfiguration, nil
}

// NewClientSocket creates a new raw client socket connected to the server
func (setup *TestServerSetup) NewClientSocket() (
	wwr.Socket,
	message.ServerConfiguration,
) {
	sock, srvConf, err := setup.ServerSetup.NewClientSocket()
	require.NoError(setup.t, err)
	return sock, srvConf
}

// CompareSessions compares a webwire session
func CompareSessions(t *testing.T, expected, actual *wwr.Session) {
	if actual == nil && expected == nil {
		return
	}

	assert.NotNil(t, expected)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Key, actual.Key)
	assert.Equal(t, expected.Creation.Unix(), actual.Creation.Unix())
}
