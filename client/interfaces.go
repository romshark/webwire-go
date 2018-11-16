package client

import (
	"context"

	webwire "github.com/qbeon/webwire-go"
)

// Client represents a webwire client instance
type Client interface {
	// Status returns the current client status
	// which is either Disabled, Disconnected or Connected.
	// The client is considered disabled when it was manually closed
	// through client.Close, while disconnected is considered
	// a temporary connection loss.
	// A disabled client won't autoconnect until enabled again
	Status() Status

	// Connect connects the client to the configured server and
	// returns an error in case of a connection failure.
	// Automatically tries to restore the previous session.
	// Enables autoconnect if it was previously disabled
	Connect() error

	// Request sends a request containing the given payload to the server and
	// asynchronously returns the servers response. It blocks until either a
	// response is received or the request fails or times out.
	//
	// Request respects context cancellation. Nil contexts are not guaranteed to
	// be accepted, context.TODO() or context.Background() should therefor
	// always be used instead.
	//
	// The returned reply object must be closed for the underlying reply message
	// buffer to be released. If the reply object is never closed it will leak
	// memory.
	Request(
		ctx context.Context,
		name []byte,
		payload webwire.Payload,
	) (webwire.Reply, error)

	// Signal sends a signal containing the given payload to the server
	Signal(
		ctx context.Context,
		name []byte,
		payload webwire.Payload,
	) error

	// Session returns an exact copy of the session object,
	// otherwise returns nil if there's currently no session
	Session() *webwire.Session

	// SessionInfo returns a copy of the session info field value
	// in the form of an empty interface to be casted to either concrete type
	SessionInfo(fieldName string) interface{}

	// PendingRequests returns the number of currently pending requests
	PendingRequests() int

	// RestoreSession tries to restore the previously opened session.
	// Fails if a session is currently already active
	RestoreSession(
		ctx context.Context,
		sessionKey []byte,
	) error

	// CloseSession disables the currently active session
	// and acknowledges the server if connected.
	// The session will be destroyed if this is it's last connection remaining.
	// If the client is not connected then the synchronization is skipped.
	// CloseSession does nothing if there's no active session
	CloseSession() error

	// Close gracefully closes the connection and disables the client.
	// A disabled client won't autoconnect until enabled again.
	Close()
}

// Implementation defines a webwire client implementation interface
type Implementation interface {
	// OnDisconnected is invoked when the client is disconnected
	// from the server for any reason.
	OnDisconnected()

	// OnSignal is invoked when the client receives a signal from the server
	OnSignal(webwire.Message)

	// OnSessionCreated is invoked when the client was assigned a new session
	OnSessionCreated(*webwire.Session)

	// OnSessionClosed is invoked when the client's session was closed
	// either by the server or the client itself
	OnSessionClosed()
}
