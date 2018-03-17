package webwire

import (
	"fmt"
	"time"
)

// ReqErrSrvShutdown represents a request error type indicating that the request cannot be
// processed due to the server currently being shut down
type ReqErrSrvShutdown struct{}

func (err ReqErrSrvShutdown) Error() string {
	return "Server is currently being shut down and won't process the request"
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
