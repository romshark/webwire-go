package webwire

import (
	"fmt"
)

// BufferOverflowErr represents a message buffer overflow error
type BufferOverflowErr struct{}

// Error implements the error interface
func (err BufferOverflowErr) Error() string {
	return "message buffer overflow"
}

// DeadlineExceededErr represents a failure due to an excess of a user-defined
// deadline
type DeadlineExceededErr struct {
	Cause error
}

// Error implements the error interface
func (err DeadlineExceededErr) Error() string {
	return err.Cause.Error()
}

// DisconnectedErr represents an error indicating a disconnected connection
type DisconnectedErr struct {
	Cause error
}

// Error implements the error interface
func (err DisconnectedErr) Error() string {
	if err.Cause == nil {
		return "disconnected"
	}
	return err.Cause.Error()
}

// DialTimeoutErr represents a dialing error caused by a timeout
type DialTimeoutErr struct{}

// Error implements the error interface
func (err DialTimeoutErr) Error() string {
	return "dial timed out"
}

// IncompatibleProtocolVersionErr represents a connection error indicating that
// the server requires an incompatible version of the protocol and can't
// therefore be connected to
type IncompatibleProtocolVersionErr struct {
	RequiredVersion  string
	SupportedVersion string
}

// Error implements the error interface
func (err IncompatibleProtocolVersionErr) Error() string {
	return fmt.Sprintf(
		"unsupported protocol version: %s (supported: %s)",
		err.RequiredVersion,
		err.SupportedVersion,
	)
}

// InternalErr represents a server-side internal error
type InternalErr struct{}

// Error implements the error interface
func (err InternalErr) Error() string {
	return "internal server error"
}

// MaxSessConnsReachedErr represents an authentication error indicating that the
// given session already reached the maximum number of concurrent connections
type MaxSessConnsReachedErr struct{}

// Error implements the error interface
func (err MaxSessConnsReachedErr) Error() string {
	return "reached maximum number of concurrent session connections"
}

// ProtocolErr represents a protocol error
type ProtocolErr struct {
	Cause error
}

// NewProtocolErr constructs a new ProtocolErr error based on the actual error
func NewProtocolErr(err error) ProtocolErr {
	return ProtocolErr{
		Cause: err,
	}
}

// Error implements the error interface
func (err ProtocolErr) Error() string {
	return err.Cause.Error()
}

// RequestErr represents an error returned when a request couldn't be processed
type RequestErr struct {
	Code    string
	Message string
}

// Error implements the error interface
func (err RequestErr) Error() string {
	return err.Message
}

// ServerShutdownErr represents a request error indicating that the request
// cannot be processed due to the server currently being shut down
type ServerShutdownErr struct{}

// Error implements the error interface
func (err ServerShutdownErr) Error() string {
	return "server is currently being shut down and won't process the request"
}

// SessionNotFoundErr represents a session restoration error indicating that the
// server didn't find the session to be restored
type SessionNotFoundErr struct{}

// Error implements the error interface
func (err SessionNotFoundErr) Error() string {
	return "session not found"
}

// SessionsDisabledErr represents an error indicating that the server has
// sessions disabled
type SessionsDisabledErr struct{}

// Error implements the error interface
func (err SessionsDisabledErr) Error() string {
	return "sessions are disabled for this server"
}

// TimeoutErr represents a failure caused by a timeout
type TimeoutErr struct {
	Cause error
}

// Error implements the error interface
func (err TimeoutErr) Error() string {
	return err.Cause.Error()
}

// TransmissionErr represents a connection error indicating a failed
// transmission
type TransmissionErr struct {
	Cause error
}

// Error implements the error interface
func (err TransmissionErr) Error() string {
	return fmt.Sprintf("message transmission failed: %s", err.Cause)
}

// CanceledErr represents a failure due to cancellation
type CanceledErr struct {
	Cause error
}

// Error implements the error interface
func (err CanceledErr) Error() string {
	return err.Cause.Error()
}

// IsTimeoutErr returns true if the given error is either a TimeoutErr
// or a DeadlineExceededErr, otherwise returns false
func IsTimeoutErr(err error) bool {
	switch err.(type) {
	case TimeoutErr:
		return true
	case DeadlineExceededErr:
		return true
	}
	return false
}

// IsCanceledErr returns true if the given error is a CanceledErr,
// otherwise returns false
func IsCanceledErr(err error) bool {
	switch err.(type) {
	case CanceledErr:
		return true
	}
	return false
}

// SockReadErr defines the interface of a webwire.Socket.Read error
type SockReadErr interface {
	error

	// IsCloseErr must return true if the error represents a closure error
	IsCloseErr() bool
}
