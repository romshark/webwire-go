package webwire

import (
	"net"
	"net/http"
	"time"
)

// SockReadErr defines the interface of a webwire.Socket.Read error
type SockReadErr interface {
	// IsAbnormalCloseErr must return true if the error represents
	// an abnormal closure error
	IsAbnormalCloseErr() bool
}

// Socket defines the abstract socket implementation interface
type Socket interface {
	// Dial must connect the socket to the specified server
	Dial(serverAddr string) error

	// Write must send the given data to the other side of the socket
	// while protecting the connection from concurrent writes
	Write(data []byte) error

	// Read must block the calling goroutine and await an incoming message.
	// When a message arrives or an error occurs Read must return
	Read() ([]byte, SockReadErr)

	// IsConnected must return true if the given socket
	// maintains an open connection or otherwise return false
	IsConnected() bool

	// RemoteAddr must return the address of the remote client
	// or nil if the client is not connected
	RemoteAddr() net.Addr

	// Close must close the socket
	Close() error

	// SetReadDeadline must set the readers deadline
	SetReadDeadline(deadline time.Time) error

	// OnPong must set the pong-message handler
	OnPong(handler func(string) error)

	// OnPing must set the ping-message handler
	OnPing(handler func(string) error)

	// WritePing must send a ping-message with the given data appended
	WritePing(data []byte, deadline time.Time) error
}

// ConnUpgrader defines the abstract interface
// of an HTTP to WebSocket connection upgrader
type ConnUpgrader interface {
	Upgrade(resp http.ResponseWriter, req *http.Request) (Socket, error)
}
