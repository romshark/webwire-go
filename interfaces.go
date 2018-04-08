package webwire

import (
	"context"
	"net/http"
)

// ServerImplementation defines the interface of a webwire server implementation
type ServerImplementation interface {
	// OnOptions is invoked when the websocket endpoint is examined by the client
	// using the HTTP OPTION method.
	OnOptions(resp http.ResponseWriter)

	// BeforeUpgrade is invoked right before the upgrade of an incoming HTTP connection request to
	// a WebSocket connection and can be used to intercept or prevent connection attempts.
	// If true is returned then the connection is normally established, though if false is returned
	// then the connection won't be established and will be canceled immediately
	BeforeUpgrade(resp http.ResponseWriter, req *http.Request) bool

	// OnClientConnected is invoked when a new client successfully established a connection
	// to the server.
	//
	// This hook will be invoked by the goroutine serving the client and thus will block the
	// initialization process, detaining the client from starting to listen for incoming messages.
	// To prevent blocking the initialization process it is advised to move any time consuming work
	// to a separate goroutine
	OnClientConnected(client *Client)

	// OnClientDisconnected is invoked when a client closes the connection to the server
	//
	// This hook will be invoked by the goroutine serving the calling client before it's suspended
	OnClientDisconnected(client *Client)

	// OnSignal is invoked when the webwire server receives a signal from a client.
	//
	// This hook will be invoked by the goroutine serving the calling client and will block any
	// other interactions with this client while executing
	OnSignal(ctx context.Context, client *Client, message *Message)

	// OnRequest is invoked when the webwire server receives a request from a client.
	// It must return either a response payload or an error.
	//
	// A webwire.ReqErr error can be returned to reply with an error code and an error message,
	// this is useful when the clients user code needs to be able to understand the error
	// and react accordingly.
	// If a non-webwire error type is returned such as an error created by fmt.Errorf(),
	// a special kind of error (internal server error) is returned to the client as a reply,
	// in this case the error will be logged and the error message will not be sent to the client
	// for security reasons as this might accidentally leak sensitive information to the client.
	//
	// This hook will be invoked by the goroutine serving the calling client and will block any
	// other interactions with this client while executing
	OnRequest(
		ctx context.Context,
		client *Client,
		message *Message,
	) (response Payload, err error)
}

// SessionManager defines the interface of a webwire server's session manager
type SessionManager interface {
	// OnSessionCreated is invoked after the synchronization of the new session
	// to the remote client.
	// The actual created session can be retrieved from the provided client agent.
	// If OnSessionCreated returns an error then this error is logged
	// but the session will not be destroyed and will remain active!
	// The only consequence of OnSessionCreation failing is that the server won't be able
	// to restore the session after the client is disconnected.
	//
	// This hook will be invoked by the goroutine calling the client.CreateSession
	// client agent method
	OnSessionCreated(client *Client) error

	// OnSessionLookup is invoked when the server is looking for a specific session given its key.
	// It must return the exact copy of the session object associated with the given key
	// for sessions to be properly restorable. If no session is found it must return nil
	// instead of the session and must not return any error.
	// If an error is returned then the it'll be logged and the session restoration will fail.
	//
	// This hook will be invoked by the goroutine serving the associated client and will block any
	// other interactions with this client while executing
	OnSessionLookup(key string) (*Session, error)

	// OnSessionClosed is invoked when the active session of the given client
	// is closed (thus destroyed) either by the server or the client through a
	// session destruction request.
	// A closed session must be permanently deleted and must not be discoverable
	// in the OnSessionLookup hook any longer.
	// If an error is returned then the it is logged.
	//
	// This hook is invoked by either a goroutine calling the client.CloseSession()
	// client agent method, or the goroutine serving the associated client, in the case of which
	// it will block any other interactions with this client while executing
	OnSessionClosed(client *Client) error
}

// SessionKeyGenerator defines the interface of a webwire servers session key generator.
// This interface must not be implemented (!) unless the default generator doesn't meet the exact
// needs of the library user, because the default generator already provides a secure implementation
type SessionKeyGenerator interface {
	// Generate is invoked when the webwire server creates a new session and requires
	// a new session key to be generated. This hook must not be used except the user
	// knows exactly what he/she does as it would compromise security if implemented improperly
	Generate() string
}

// SessionInfo represents a session info object implementation interface.
// It defines a set of important methods that must be implemented carefully
// in order to avoid race conditions
type SessionInfo interface {
	// Fields must return the exact names of all fields
	// of the session info object. This getter method must be idempotent,
	// which means that it must always return the same list of names
	Fields() []string

	// Value must return an exact deep copy of the value of a session info
	// object field identified by the given field name.
	//
	// Note that returning a shallow copy (such as shallow copies of
	// maps or slices for example) could lead to potentially dangerous
	// race conditions and undefined behavior
	Value(fieldName string) interface{}

	// Copy must return an exact deep copy of the entire session info object.
	//
	// Note that returning a shallow copy (such as shallow copies of
	// maps or slices for example) could lead to potentially dangerous
	// race conditions and undefined behavior
	Copy() SessionInfo
}
