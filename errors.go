package webwire

import (
	"github.com/qbeon/webwire-go/wwrerr"
)

// CanceledErr represents a failure due to cancellation
type CanceledErr = wwrerr.CanceledErr

// DeadlineExceededErr represents a failure due to an excess of a user-defined
// deadline
type DeadlineExceededErr = wwrerr.DeadlineExceededErr

// DisconnectedErr represents an error indicating a disconnected connection
type DisconnectedErr = wwrerr.DisconnectedErr

// DialTimeoutErr represents a dialing error caused by a timeout
type DialTimeoutErr = wwrerr.DialTimeoutErr

// IncompatibleProtocolVersionErr represents a connection error indicating that
// the server requires an incompatible version of the protocol and can't
// therefore be connected to
type IncompatibleProtocolVersionErr = wwrerr.IncompatibleProtocolVersionErr

// InternalErr represents a server-side internal error
type InternalErr = wwrerr.InternalErr

// MaxSessConnsReachedErr represents an authentication error indicating that the
// given session already reached the maximum number of concurrent connections
type MaxSessConnsReachedErr = wwrerr.MaxSessConnsReachedErr

// ProtocolErr represents a protocol error
type ProtocolErr = wwrerr.ProtocolErr

// RequestErr represents an error returned when a request couldn't be processed
type RequestErr = wwrerr.RequestErr

// ServerShutdownErr represents a request error indicating that the request
// cannot be processed due to the server currently being shut down
type ServerShutdownErr = wwrerr.ServerShutdownErr

// SessionNotFoundErr represents a session restoration error indicating that the
// server didn't find the session to be restored
type SessionNotFoundErr = wwrerr.SessionNotFoundErr

// SessionsDisabledErr represents an error indicating that the server has
// sessions disabled
type SessionsDisabledErr = wwrerr.SessionsDisabledErr

// TimeoutErr represents a failure caused by a timeout
type TimeoutErr = wwrerr.TimeoutErr

// TransmissionErr represents a connection error indicating a failed
// transmission
type TransmissionErr = wwrerr.TransmissionErr

// BufferOverflowErr represents a message buffer overflow error
type BufferOverflowErr = wwrerr.BufferOverflowErr

// IsTimeoutErr returns true if the given error is either a TimeoutErr
// or a DeadlineExceededErr, otherwise returns false
func IsTimeoutErr(err error) bool { return wwrerr.IsTimeoutErr(err) }

// IsCanceledErr returns true if the given error is a CanceledErr,
// otherwise returns false
func IsCanceledErr(err error) bool { return wwrerr.IsCanceledErr(err) }
