package webwire

import (
	"fmt"
)

// ErrBufferOverflow represents a message buffer overflow error
type ErrBufferOverflow struct{}

// Error implements the error interface
func (err ErrBufferOverflow) Error() string {
	return "message buffer overflow"
}

// ErrDeadlineExceeded represents a failure due to an excess of a user-defined
// deadline
type ErrDeadlineExceeded struct {
	Cause error
}

// Error implements the error interface
func (err ErrDeadlineExceeded) Error() string {
	return err.Cause.Error()
}

// ErrDisconnected represents an error indicating a disconnected connection
type ErrDisconnected struct {
	Cause error
}

// Error implements the error interface
func (err ErrDisconnected) Error() string {
	if err.Cause == nil {
		return "disconnected"
	}
	return err.Cause.Error()
}

// ErrDialTimeout represents a dialing error caused by a timeout
type ErrDialTimeout struct{}

// Error implements the error interface
func (err ErrDialTimeout) Error() string {
	return "dial timed out"
}

// ErrIncompatibleProtocolVersion represents a connection error indicating that
// the server requires an incompatible version of the protocol and can't
// therefore be connected to
type ErrIncompatibleProtocolVersion struct {
	RequiredVersion  string
	SupportedVersion string
}

// Error implements the error interface
func (err ErrIncompatibleProtocolVersion) Error() string {
	return fmt.Sprintf(
		"unsupported protocol version: %s (supported: %s)",
		err.RequiredVersion,
		err.SupportedVersion,
	)
}

// ErrInternal represents a server-side internal error
type ErrInternal struct{}

// Error implements the error interface
func (err ErrInternal) Error() string {
	return "internal server error"
}

// ErrMaxSessConnsReached represents an authentication error indicating that the
// given session already reached the maximum number of concurrent connections
type ErrMaxSessConnsReached struct{}

// Error implements the error interface
func (err ErrMaxSessConnsReached) Error() string {
	return "reached maximum number of concurrent session connections"
}

// ErrProtocol represents a protocol error
type ErrProtocol struct {
	Cause error
}

// ErrNewProtocol constructs a new ErrProtocol error based on the actual error
func ErrNewProtocol(err error) ErrProtocol {
	return ErrProtocol{
		Cause: err,
	}
}

// Error implements the error interface
func (err ErrProtocol) Error() string {
	return err.Cause.Error()
}

// ErrRequest represents an error returned when a request couldn't be processed
type ErrRequest struct {
	Code    string
	Message string
}

// Error implements the error interface
func (err ErrRequest) Error() string {
	return err.Message
}

// ErrServerShutdown represents a request error indicating that the request
// cannot be processed due to the server currently being shut down
type ErrServerShutdown struct{}

// Error implements the error interface
func (err ErrServerShutdown) Error() string {
	return "server is currently being shut down and won't process the request"
}

// ErrSessionNotFound represents a session restoration error indicating that the
// server didn't find the session to be restored
type ErrSessionNotFound struct{}

// Error implements the error interface
func (err ErrSessionNotFound) Error() string {
	return "session not found"
}

// ErrSessionsDisabled represents an error indicating that the server has
// sessions disabled
type ErrSessionsDisabled struct{}

// Error implements the error interface
func (err ErrSessionsDisabled) Error() string {
	return "sessions are disabled for this server"
}

// ErrTimeout represents a failure caused by a timeout
type ErrTimeout struct {
	Cause error
}

// Error implements the error interface
func (err ErrTimeout) Error() string {
	return err.Cause.Error()
}

// ErrTransmission represents a connection error indicating a failed
// transmission
type ErrTransmission struct {
	Cause error
}

// Error implements the error interface
func (err ErrTransmission) Error() string {
	return fmt.Sprintf("message transmission failed: %s", err.Cause)
}

// ErrCanceled represents a failure due to cancellation
type ErrCanceled struct {
	Cause error
}

// Error implements the error interface
func (err ErrCanceled) Error() string {
	return err.Cause.Error()
}

// IsErrTimeout returns true if the given error is either a ErrTimeout
// or a ErrDeadlineExceeded, otherwise returns false
func IsErrTimeout(err error) bool {
	switch err.(type) {
	case ErrTimeout:
		return true
	case ErrDeadlineExceeded:
		return true
	}
	return false
}

// IsErrCanceled returns true if the given error is a ErrCanceled,
// otherwise returns false
func IsErrCanceled(err error) bool {
	switch err.(type) {
	case ErrCanceled:
		return true
	}
	return false
}

// ErrSockRead defines the interface of a webwire.Socket.Read error
type ErrSockRead interface {
	error

	// IsCloseErr must return true if the error represents a closure error
	IsCloseErr() bool
}
