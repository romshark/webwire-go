package webwire

import (
	"net"
)

// Server defines the interface of a webwire server instance
type Server interface {

	// Run will luanch the webwire server blocking the calling goroutine
	// until the server is either gracefully shut down
	// or crashes returning an error
	Run() error

	// Addr returns the address the webwire server is listening on
	Addr() net.Addr

	// Shutdown appoints a server shutdown and blocks the calling goroutine
	// until the server is gracefully stopped awaiting all currently processed
	// signal and request handlers to return.
	// During the shutdown incoming connections are rejected
	// with 503 service unavailable.
	// Incoming requests are rejected with an error while incoming signals
	// are just ignored
	Shutdown() error

	// ActiveSessionsNum returns the number of currently active sessions
	ActiveSessionsNum() int

	// SessionConnectionsNum implements the SessionRegistry interface
	SessionConnectionsNum(sessionKey string) int

	// SessionConnections implements the SessionRegistry interface
	SessionConnections(sessionKey string) []*Client

	// CloseSession closes the session identified by the given key
	// and returns the number of closed connections.
	// If there was no session found -1 is returned
	CloseSession(sessionKey string) int
}
