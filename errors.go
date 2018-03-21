package webwire

import (
	"fmt"
	"time"
)

// ConnIncompErr represents a connection error type indicating that the server
// requires an incompatible version of the protocol and can't therefore be connected to.
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

// NewConnIncompErr constructs and returns a new incompatible protocol version error
// based on the required and supported protocol versions
func NewConnIncompErr(requiredVersion, supportedVersion string) ConnIncompErr {
	return ConnIncompErr{
		requiredVersion:  requiredVersion,
		supportedVersion: supportedVersion,
	}
}

// ConnDialErr represents a connection error type indicating that the dialing failed.
type ConnDialErr struct {
	msg string
}

func (err ConnDialErr) Error() string {
	return err.msg
}

// NewConnDialErr constructs and returns a new connection dial error
// based on the actual error message
func NewConnDialErr(err error) ConnDialErr {
	return ConnDialErr{
		msg: err.Error(),
	}
}

// ReqTransErr represents a connection error type indicating that the dialing failed.
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

// ReqSrvShutdownErr represents a request error type indicating that the request cannot be
// processed due to the server currently being shut down
type ReqSrvShutdownErr struct{}

func (err ReqSrvShutdownErr) Error() string {
	return "Server is currently being shut down and won't process the request"
}

// ReqInternalErr represents a request error type indicating that the request failed due
// to an internal server-side error
type ReqInternalErr struct{}

func (err ReqInternalErr) Error() string {
	return "Internal server error"
}

// ReqTimeoutErr represents a request error type indicating that the server
// wasn't able to reply within the given time frame causing the request to time out.
type ReqTimeoutErr struct {
	Target time.Duration
}

func (err ReqTimeoutErr) Error() string {
	return fmt.Sprintf("Server didn't manage to reply within %s", err.Target)
}

// ReqErr represents an error returned in case of a request that couldn't be processed
type ReqErr struct {
	Code    string `json:"c"`
	Message string `json:"m,omitempty"`
}

func (err ReqErr) Error() string {
	return err.Message
}

// SessionsDisabledErr represents an error type indicating that the server has sessions disabled
type SessionsDisabledErr struct{}

func (err SessionsDisabledErr) Error() string {
	return "Sessions are disabled for this server"
}

// SessNotFoundErr represents a session restoration error type indicating that the server didn't
// find the session to be restored
type SessNotFoundErr struct{}

func (err SessNotFoundErr) Error() string {
	return "Session not found"
}

// MaxSessConnsReachedErr represents an authentication error type indicating that the given session
// already reached the maximum number of concurrent connections
type MaxSessConnsReachedErr struct{}

func (err MaxSessConnsReachedErr) Error() string {
	return "Reached maximum number of concurrent session connections"
}

// DisconnectedErr represents an error type indicating that the targeted client is disconnected
type DisconnectedErr struct {
	msg string
}

func (err DisconnectedErr) Error() string {
	return err.msg
}
