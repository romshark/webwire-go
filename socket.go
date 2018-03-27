package webwire

import (
	"net"
	"net/http"
)

// SockReadErr defines the interface of a webwire.Socket.Read error
type SockReadErr interface {
	// IsAbnormalCloseErr must return true if the error represents an abnormal closure error
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
	// When a message arrives or an error occurs Read must return.
	Read() ([]byte, SockReadErr)

	// IsConnected must return true if the given socket maintains an open connection,
	// otherwise return false
	IsConnected() bool

	// RemoteAddr must return the address of the remote client
	// or nil if the client is not connected
	RemoteAddr() net.Addr

	// Close must close the socket
	Close() error
}

// ConnUpgrader defines the abstract interface of an HTTP to WebSocket connection upgrader
type ConnUpgrader interface {
	Upgrade(resp http.ResponseWriter, req *http.Request) (Socket, error)
}
