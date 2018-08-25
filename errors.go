package webwire

import (
	"fmt"
)

// ConnIncompErr represents a connection error type indicating that the server
// requires an incompatible version of the protocol
// and can't therefore be connected to.
type ConnIncompErr struct {
	requiredVersion  string
	supportedVersion string
}

func (err ConnIncompErr) Error() string {
	return fmt.Sprintf(
		"Unsupported protocol version: %s (%s is supported by this client)",
		err.requiredVersion,
		err.supportedVersion,
	)
}

// NewConnIncompErr constructs and returns a new "incompatible protocol version"
// error based on the required and supported protocol versions
func NewConnIncompErr(requiredVersion, supportedVersion string) ConnIncompErr {
	return ConnIncompErr{
		requiredVersion:  requiredVersion,
		supportedVersion: supportedVersion,
	}
}

// ReqTransErr represents a connection error type
// indicating that the dialing failed.
type ReqTransErr struct {
	msg string
}

func (err ReqTransErr) Error() string {
	return fmt.Sprintf("Message transmission failed: %s", err.msg)
}

// NewReqTransErr constructs and returns a new request transmission error
// based on the actual error message
func NewReqTransErr(err error) ReqTransErr {
	return ReqTransErr{
		msg: err.Error(),
	}
}

// ReqSrvShutdownErr represents a request error type indicating
// that the request cannot be processed
// due to the server currently being shut down
type ReqSrvShutdownErr struct{}

func (err ReqSrvShutdownErr) Error() string {
	return "Server is currently being shut down and won't process the request"
}

// ReqInternalErr represents a request error type
// indicating that the request failed due to an internal server-side error
type ReqInternalErr struct{}

func (err ReqInternalErr) Error() string {
	return "Internal server error"
}

// TimeoutErr represents a failure due to a timeout
type TimeoutErr struct {
	cause error
}

// NewTimeoutErr constructs a new timeout error
func NewTimeoutErr(cause error) TimeoutErr {
	return TimeoutErr{cause}
}

// Error implements the error interface
func (err TimeoutErr) Error() string {
	return err.cause.Error()
}

// DeadlineExceededErr represents a failure due to
// an excess of a user-defined deadline
type DeadlineExceededErr struct {
	cause error
}

// NewDeadlineExceededErr constructs a new deadline excess error
func NewDeadlineExceededErr(cause error) DeadlineExceededErr {
	return DeadlineExceededErr{cause}
}

// Error implements the error interface
func (err DeadlineExceededErr) Error() string {
	return err.cause.Error()
}

// ReqErr represents an error returned in case of
// a request that couldn't be processed
type ReqErr struct {
	Code    string
	Message string
}

func (err ReqErr) Error() string {
	return err.Message
}

// SessionsDisabledErr represents an error type
// indicating that the server has sessions disabled
type SessionsDisabledErr struct{}

func (err SessionsDisabledErr) Error() string {
	return "Sessions are disabled for this server"
}

// SessNotFoundErr represents a session restoration error type
// indicating that the server didn't find the session to be restored
type SessNotFoundErr struct{}

func (err SessNotFoundErr) Error() string {
	return "Session not found"
}

// MaxSessConnsReachedErr represents an authentication error type
// indicating that the given session already reached the maximum number
// of concurrent connections
type MaxSessConnsReachedErr struct{}

func (err MaxSessConnsReachedErr) Error() string {
	return "Reached maximum number of concurrent session connections"
}

// DisconnectedErr represents an error type
// indicating that the targeted client is disconnected
type DisconnectedErr struct {
	Cause error
}

// NewDisconnectedErr constructs a new DisconnectedErr error
// based on the actual error
func NewDisconnectedErr(err error) DisconnectedErr {
	return DisconnectedErr{
		Cause: err,
	}
}

func (err DisconnectedErr) Error() string {
	if err.Cause == nil {
		return "Disconnected"
	}
	return err.Cause.Error()
}

// ProtocolErr represents an error type
// indicating an error in the protocol implementation
type ProtocolErr struct {
	cause error
}

// NewProtocolErr constructs a new ProtocolErr error based on the actual error
func NewProtocolErr(err error) ProtocolErr {
	return ProtocolErr{
		cause: err,
	}
}

func (err ProtocolErr) Error() string {
	return err.cause.Error()
}

// CanceledErr represents a failure due to cancelation
type CanceledErr struct {
	cause error
}

// NewCanceledErr constructs a new canceled errror instance
func NewCanceledErr(cause error) CanceledErr {
	return CanceledErr{cause}
}

// Error implements the error interface
func (err CanceledErr) Error() string {
	return err.cause.Error()
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
