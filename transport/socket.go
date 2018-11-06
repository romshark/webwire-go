package transport

import (
	"io"
	"net"
	"net/url"
	"time"

	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/wwrerr"
)

// Socket defines the abstract socket implementation interface
type Socket interface {
	// Dial connects the socket to the specified server
	Dial(serverAddr url.URL) error

	// GetWriter returns a writer for the next message to send. The writer's
	// Close method flushes the complete message to the network
	GetWriter() (io.WriteCloser, error)

	// Read blocks the calling goroutine and awaits an incoming message
	Read(*message.Message) wwrerr.SockReadErr

	// IsConnected returns true if the given socket maintains an open
	// connection or otherwise return false
	IsConnected() bool

	// RemoteAddr returns the address of the remote client or nil if the client
	// is not connected
	RemoteAddr() net.Addr

	// SetReadDeadline sets the readers deadline
	SetReadDeadline(deadline time.Time) error

	// Close closes the socket
	Close() error
}
