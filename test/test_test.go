package test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/fasthttp/websocket"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	wwrtrn "github.com/qbeon/webwire-go/transport"
	wwrfasthttp "github.com/qbeon/webwire-go/transport/fasthttp"
	wwrmemchan "github.com/qbeon/webwire-go/transport/memchan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

type serverSetup struct {
	Transport wwrtrn.Transport
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
	trans wwrtrn.Transport,
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
	switch *argTransport {
	case "fasthttp/websocket":
		if trans == nil {
			// Use default configuration
			trans = &wwrfasthttp.Transport{
				HTTPServer: &fasthttp.Server{
					ReadBufferSize:  1024 * 8,
					WriteBufferSize: 1024 * 8,
				},
			}
		} else {
			if _, isType := trans.(*wwrfasthttp.Transport); !isType {
				return serverSetup{}, fmt.Errorf(
					"unexpected server transport implementation: %s",
					reflect.TypeOf(trans),
				)
			}
		}
	case "memchan":
		if trans == nil {
			// Use default configuration
			trans = &wwrmemchan.Transport{}
		} else {
			if _, isType := trans.(*wwrmemchan.Transport); !isType {
				return serverSetup{}, fmt.Errorf(
					"unexpected server transport implementation: %s",
					reflect.TypeOf(trans),
				)
			}
		}
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
	trans wwrtrn.Transport,
) testServerSetup {
	setup, err := setupServer(impl, opts, trans)
	require.NoError(t, err)
	return testServerSetup{t, setup}
}

// newClient sets up a new test client instance
func (setup *serverSetup) newClient(
	options wwrclt.Options,
	transport wwrtrn.ClientTransport,
	hooks testClientHooks,
) (*testClient, error) {
	return newClient(
		setup.Server.Address(),
		setup.Transport,
		options,
		transport,
		hooks,
	)
}

// newClient sets up a new test client instance
func (setup *testServerSetup) newClient(
	options wwrclt.Options,
	transport wwrtrn.ClientTransport,
	hooks testClientHooks,
) *testClient {
	clt, err := newClient(
		setup.Server.Address(),
		setup.Transport,
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
	transport wwrtrn.Transport,
	options wwrclt.Options,
	clientTransport wwrtrn.ClientTransport,
	hooks testClientHooks,
) (*testClient, error) {
	// Prepare transport layer
	switch transport.(type) {
	case *wwrfasthttp.Transport:
		if clientTransport == nil {
			// Use default configuration
			clientTransport = &wwrfasthttp.ClientTransport{}
		} else {
			_, isType := clientTransport.(*wwrfasthttp.ClientTransport)
			if !isType {
				return nil, fmt.Errorf(
					"unexpected client transport implementation: %s",
					reflect.TypeOf(clientTransport),
				)
			}
		}

	case *wwrmemchan.Transport:
		if clientTransport == nil {
			// Use default configuration
			clientTransport = &wwrmemchan.ClientTransport{
				Server: transport.(*wwrmemchan.Transport),
			}
		} else {
			_, isType := clientTransport.(*wwrmemchan.ClientTransport)
			if !isType {
				return nil, fmt.Errorf(
					"unexpected client transport implementation: %s",
					reflect.TypeOf(clientTransport),
				)
			}
		}
		// Rewrite server reference
		clientTransport.(*wwrmemchan.ClientTransport).Server =
			transport.(*wwrmemchan.Transport)

	default:
		return nil, fmt.Errorf(
			"unexpected server transport implementation: %s",
			reflect.TypeOf(transport),
		)
	}

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
func (setup *serverSetup) newClientSocket() (wwrtrn.Socket, error) {
	switch srvTrans := setup.Transport.(type) {
	case *wwrfasthttp.Transport:
		// Setup a regular websocket connection
		serverAddr := setup.Server.Address()
		if serverAddr.Scheme == "https" {
			serverAddr.Scheme = "wss"
		} else {
			serverAddr.Scheme = "ws"
		}

		conn, _, err := websocket.DefaultDialer.Dial(serverAddr.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("dialing failed: %s", err)
		}

		return wwrfasthttp.NewConnectedSocket(conn), nil

	case *wwrmemchan.Transport:
		_, sock := wwrmemchan.NewEntangledSockets(srvTrans)
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
func (setup *testServerSetup) newClientSocket() wwrtrn.Socket {
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

var argTransport = flag.String(
	"wwr-transport",
	"fasthttp/websocket",
	"determines the webwire transport layer implementation",
)

// parseArgs parses and validates the CLI argument
func parseArgs() {
	flag.Parse()

	switch *argTransport {
	case "fasthttp/websocket":
	case "memchan":
	default:
		log.Fatalf(
			"unknown transport layer implementation: '%s'",
			*argTransport,
		)
	}
}

// TestMain executes the tests
func TestMain(m *testing.M) {
	parseArgs()

	fmt.Printf("transport: %s\n", *argTransport)

	os.Exit(m.Run())
}
