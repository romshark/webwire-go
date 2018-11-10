package transport

import (
	"io"
	"net"
	"time"

	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/wwrerr"
)

// Socket defines the abstract socket implementation interface
type Socket interface {
	// GetWriter returns a writer for the next message to send. The writer's
	// Close method flushes the written message to the network. In case of
	// concurrent use GetWriter will block until the previous writer is closed
	// and a new one is available
	GetWriter() (io.WriteCloser, error)

	// Read blocks the calling goroutine and awaits an incoming message. If
	// deadline is 0 then Read will never timeout. In case of concurrent use
	// Read will block until the previous call finished
	Read(into *message.Message, deadline time.Time) wwrerr.SockReadErr

	// IsConnected returns true if the given socket maintains an open connection
	// or otherwise return false
	IsConnected() bool

	// RemoteAddr returns the address of the remote client or nil if the client
	// is not connected
	RemoteAddr() net.Addr

	// Close closes the socket
	Close() error
}
