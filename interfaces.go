package webwire

import (
	"context"
	"net"
	"net/url"
	"time"

	"github.com/qbeon/webwire-go/connopt"
)

// ConnectionAcceptance defines whether a connection is to be accepted
type ConnectionAcceptance = connopt.ConnectionAcceptance

// Accept instructs the server to accept the incoming connection
const Accept = connopt.Accept

// Refuse instructs the server to refuse the incoming connection
const Refuse = connopt.Refuse

// ConnectionOptions represents the options applied to an individual connection
// during accept
type ConnectionOptions = connopt.ConnectionOptions

// HeadlessServer defines the interface of a headless webwire server instance
type HeadlessServer interface {
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
	SessionConnections(sessionKey string) []Connection

	// CloseSession closes the session identified by the given key and returns
	// the affected connections, a list of errors for each session session
	// closure attempt and a general error which is not nil if at least
	// of the closeErrors errors is not nil.
	// If no session was closed then (nil, nil, nil) is returned.
	CloseSession(sessionKey string) (
		affectedConnections []Connection,
		closeErrors []error,
		err error,
	)
}

// Server defines the interface of a headed webwire server instance
type Server interface {
	// Run will launch the webwire server blocking the calling goroutine
	// until the server is either gracefully shut down
	// or crashes returning an error
	Run() error

	// Address returns the URL address the webwire server is listening on
	Address() url.URL

	HeadlessServer
}

// ServerImplementation defines the interface
// of a webwire server implementation
type ServerImplementation interface {
	// OnClientConnected is invoked when a new client successfully established
	// a connection to the server.
	//
	// This hook will be invoked by the goroutine serving the client and thus
	// will block the initialization process, detaining the client from
	// starting to listen for incoming messages.
	// To prevent blocking the initialization process it is advised to move
	// any time consuming work to a separate goroutine
	OnClientConnected(connectionOptions ConnectionOptions, client Connection)

	// OnClientDisconnected is invoked when a client closes the connection
	// to the server.
	//
	// This hook will be invoked by the goroutine serving
	// the calling client before it's suspended
	OnClientDisconnected(client Connection, reason error)

	// OnSignal is invoked when the webwire server receives
	// a signal from a client.
	//
	// This hook will be invoked by the goroutine serving the calling client
	// and will block any other interactions with this client while executing
	OnSignal(ctx context.Context, client Connection, message Message)

	// OnRequest is invoked when the webwire server receives a request
	// from a client. It must return either a response payload or an error.
	//
	// A webwire.ReqErr error can be returned to reply with an error code
	// and an error message, this is useful when the clients user code needs
	// to be able to understand the error and react accordingly.
	// If a non-webwire error type is returned such as an error created by
	// fmt.Errorf() then a special kind of error (internal server error)
	// is returned to the client as a reply, in this case the error will be
	// logged and the error message will not be sent to the client
	// for security reasons as this might accidentally leak
	// sensitive information to the client.
	//
	// This hook will be invoked by the goroutine serving the calling client
	// and will block any other interactions with this client while executing
	OnRequest(
		ctx context.Context,
		client Connection,
		message Message,
	) (
		payload Payload,
		err error,
	)
}

// Connection represents a connected client
type Connection interface {
	// IsActive returns true if this connection is in active state
	// ready to accept incoming messages, otherwise returns false
	IsActive() bool

	// RemoteAddr returns the remote address
	RemoteAddr() net.Addr

	// Creation returns the time of connection establishment
	Creation() time.Time

	// Info returns arbitrary information by key which was assigned by the
	// transport layer implementation during the connection establishment.
	// Returns nil if the provided key was not found
	Info(key int) interface{}

	// Signal sends a named signal containing the given payload to the client
	Signal(name []byte, payload Payload) error

	// CreateSession creates a new session for this connection and
	// automatically synchronizes the new session to the remote client.
	// The synchronization happens asynchronously using a signal
	// and doesn't block the calling goroutine.
	// Returns an error if there's already another session active
	CreateSession(attachment SessionInfo) error

	// CloseSession disables the currently active session for this connection
	// and synchronize the closure to the remote client.
	// The session will be destroyed if this is it's last connection remaining.
	// Does nothing if there's no active session
	CloseSession() error

	// HasSession returns true if this connection currently has
	// a session assigned, otherwise returns false
	HasSession() bool

	// Session returns an exact copy of the session object or nil if there's no
	// session currently assigned to this connection
	Session() *Session

	// SessionKey returns the key of the currently assigned session.
	// Returns an empty string if there's no session assigned to this connection
	SessionKey() string

	// SessionCreation returns the creation time
	// of the currently assigned session.
	// Warning: be sure to check whether there's a session beforehand because
	// this function will return garbage if
	// there's currently no session assigned
	SessionCreation() time.Time

	// SessionInfo returns a copy of the session info field value
	// in the form of an empty interface to be casted to either concrete type
	SessionInfo(name string) interface{}

	// Close marks this connection for shutdown.
	// It defers closing the connection until all work on it is done
	// and removes it from the session registry.
	// Does nothing when called multiple times
	Close()
}

// SessionLookupResult represents the result of a session lookup
type SessionLookupResult interface {
	// Creation returns the retrieved creation time
	Creation() time.Time

	// LastLookup returns the retrieved last lookup time
	LastLookup() time.Time

	// Info returns the retrieved session information
	// in the form of a JSON compliant variant map
	Info() map[string]interface{}
}

// SessionManager defines the interface of a webwire server's session manager
type SessionManager interface {
	// OnSessionCreated is invoked after the synchronization of the new session
	// to the remote client.
	// The actual created session can be retrieved from the provided connection.
	// If OnSessionCreated returns an error then this error is logged
	// but the session will not be destroyed and will remain active!
	// The only consequence of OnSessionCreation failing is that the server
	// won't be able to restore the session after the client is disconnected.
	//
	// This hook will be invoked by the goroutine calling the
	// client.CreateSession connection method
	OnSessionCreated(client Connection) error

	// OnSessionLookup is invoked when the server is looking for a specific
	// session given its key. If the session wasn't found then nil will be
	// returned without an error, otherwise if the session was found then a
	// SessionLookupResult instance created by wwr.NewSessionLookupResult will
	// be returned with the required parameters set.
	//
	// If an error (that's not a webwire.SessNotFoundErr) is returned then it'll
	// be logged and the session restoration will fail.
	//
	// This hook will be invoked by the goroutine serving the associated client
	// and will block any other interactions with this client while executing
	//
	// WARNING: if this hooks doesn't update the LastLookup field of the found
	// session object then the session garbage collection won't work properly
	OnSessionLookup(key string) (result SessionLookupResult, err error)

	// OnSessionClosed is invoked when the session associated with the given key
	// is closed (thus destroyed) either by the server or the client. A closed
	// session is permanently deleted and becomes undiscoverable in the
	// OnSessionLookup hook. A returned error is logged to the wwr error log.
	//
	// This hook is invoked by either a goroutine calling the method
	// connection.CloseSession(), or the goroutine serving the associated
	// client, in the case of which it will block any other interactions with
	// this client while executing
	OnSessionClosed(sessionKey string) error
}

// SessionKeyGenerator defines the interface of a webwire server's
// session key generator. This interface must not be implemented (!) unless
// the default generator doesn't meet the exact needs of the library user,
// because the default generator already provides a secure implementation
type SessionKeyGenerator interface {
	// Generate is invoked when the webwire server creates a new session
	// and requires a new session key to be generated.
	// This hook must not be used except the user knows exactly what he/she does
	// as it would compromise security if implemented improperly
	Generate() string
}

// SessionInfo represents a session info object implementation interface.
// It defines a set of important methods that must be implemented carefully
// in order to avoid race conditions
type SessionInfo interface {
	// Fields returns the exact names of all fields of the session info object.
	// This getter method is idempotent, which means that it always returns the
	// same list of names
	Fields() []string

	// Value returns an exact deep copy of the value of a session info object
	// field identified by the given field name.
	//
	// Note that returning a shallow copy (such as shallow copies of maps or
	// slices for example) could lead to potentially dangerous race conditions
	// and undefined behavior
	Value(fieldName string) interface{}

	// Copy returns an exact deep copy of the entire session info object.
	//
	// Note that returning a shallow copy (such as shallow copies of
	// maps or slices for example) could lead to potentially dangerous
	// race conditions and undefined behavior
	Copy() SessionInfo
}

// SessionInfoParser represents the type of a session info parser function. The
// session info parser is invoked during the parsing of a newly assigned session
// on the client, as well as during the parsing of a saved serialized session.
// It returns a webwire.SessionInfo compliant object constructed from the data
// given
type SessionInfoParser func(map[string]interface{}) SessionInfo

// Message represents a WebWire message
//
// Message is not thread-safe and must not be used concurrently from within
// multiple goroutines
type Message interface {
	// Identifier returns the message identifier
	Identifier() [8]byte

	// Name returns the name of the message
	Name() []byte

	// PayloadEncoding returns the payload encoding type
	PayloadEncoding() PayloadEncoding

	// Payload returns the message payload in binary format
	Payload() []byte

	// PayloadUtf8 returns the message payload in textual UTF8 format
	PayloadUtf8() ([]byte, error)

	// Close closes the message releasing the underlying buffer
	Close()
}

// Reply defines a webwire request reply message
//
// Reply is not thread-safe and must not be used concurrently from within
// multiple goroutines
type Reply interface {
	// PayloadEncoding returns the payload encoding type
	PayloadEncoding() PayloadEncoding

	// Payload returns the message payload in binary format
	Payload() []byte

	// PayloadUtf8 returns the message payload in textual UTF8 format
	PayloadUtf8() ([]byte, error)

	// Close closes the reply message releasing the underlying buffer
	Close()
}
