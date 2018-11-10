package transport

import (
	"net/url"
	"time"

	"github.com/qbeon/webwire-go/connopt"
)

// IsShuttingDown must be called when the server is accepting a new connection
// and refuse the connection if true is returned
type IsShuttingDown func() bool

// OnNewConnection must be called when the connection is ready to be used by the
// webwire server
type OnNewConnection func(
	connectionOptions connopt.ConnectionOptions,
	clientAgentString []byte,
	socket Socket,
)

// Transport defines the interface of a webwire transport
type Transport interface {
	// Initialize initializes the server
	Initialize(
		host string,
		readTimeout time.Duration,
		messageBufferSize uint32,
		isShuttingdown IsShuttingDown,
		onNewConnection OnNewConnection,
	) error

	// Serve starts serving blocking the calling goroutine
	Serve() error

	// Shutdown shuts the server down
	Shutdown() error

	// Address returns the URL address the server is listening on
	Address() url.URL
}

// ClientTransport defines the interface of a webwire client transport
type ClientTransport interface {
	// NewSocket initializes a new client socket
	NewSocket(dialTimeout time.Duration) (ClientSocket, error)
}
