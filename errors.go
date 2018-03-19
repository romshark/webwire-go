package webwire

import (
	"fmt"
	"time"
)

// ConnErrIncomp represents a connection error type indicating that the server
// requires an incompatible version of the protocol and can't therefore be connected to.
type ConnErrIncomp struct {
	requiredVersion  string
	supportedVersion string
}

func (err ConnErrIncomp) Error() string {
	return fmt.Sprintf(
		"Unsupported protocol version: %s (%s is supported by this client)",
		err.requiredVersion,
		err.supportedVersion,
	)
}

// NewConnErrIncomp constructs and returns a new incompatible protocol version error
// based on the required and supported protocol versions
func NewConnErrIncomp(requiredVersion, supportedVersion string) ConnErrIncomp {
	return ConnErrIncomp{
		requiredVersion:  requiredVersion,
		supportedVersion: supportedVersion,
	}
}

// ConnErrDial represents a connection error type indicating that the dialing failed.
type ConnErrDial struct {
	msg string
}

func (err ConnErrDial) Error() string {
	return err.msg
}

// NewConnErrDial constructs and returns a new connection dial error
// based on the actual error message
func NewConnErrDial(err error) ConnErrDial {
	return ConnErrDial{
		msg: err.Error(),
	}
}

// ReqErrTrans represents a connection error type indicating that the dialing failed.
type ReqErrTrans struct {
	msg string
}

func (err ReqErrTrans) Error() string {
	return fmt.Sprintf("Message transmission failed: %s", err.msg)
}

// NewReqErrTrans constructs and returns a new request transmission error
// based on the actual error message
func NewReqErrTrans(err error) ReqErrTrans {
	return ReqErrTrans{
		msg: err.Error(),
	}
}

// ReqErrSrvShutdown represents a request error type indicating that the request cannot be
// processed due to the server currently being shut down
type ReqErrSrvShutdown struct{}

func (err ReqErrSrvShutdown) Error() string {
	return "Server is currently being shut down and won't process the request"
}

// ReqErrInternal represents a request error type indicating that the request failed due
// to an internal server-side error
type ReqErrInternal struct{}

func (err ReqErrInternal) Error() string {
	return "Internal server error"
}

// ReqErrTimeout represents a request error type indicating that the server
// wasn't able to reply within the given time frame causing the request to time out.
type ReqErrTimeout struct {
	Target time.Duration
}

func (err ReqErrTimeout) Error() string {
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

// SessionsDisabled represents an error type indicating that the server has sessions disabled
type SessionsDisabled struct{}

func (err SessionsDisabled) Error() string {
	return "Sessions are disabled for this server"
}

// SessNotFound represents a session restoration error type indicating that the server didn't
// find the session to be restored
type SessNotFound struct{}

func (err SessNotFound) Error() string {
	return "Session not found"
}

// MaxSessConnsReached represents an authentication error type indicating that the given session
// already reached the maximum number of concurrent connections
type MaxSessConnsReached struct{}

func (err MaxSessConnsReached) Error() string {
	return "Reached maximum number of concurrent session connections"
}
