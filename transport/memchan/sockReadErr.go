package memchan

import "fmt"

// ErrSockRead implements the ErrSockRead interface
type ErrSockRead struct {
	// closed is true when the error was caused by a graceful socket closure
	closed bool

	err error
}

// Error implements the Go error interface
func (err ErrSockRead) Error() string {
	if err.closed {
		return "socket closed"
	}
	return fmt.Sprintf("reading socket failed: %s", err.err)
}

// IsCloseErr implements the ErrSockRead interface
func (err ErrSockRead) IsCloseErr() bool {
	return err.closed
}
