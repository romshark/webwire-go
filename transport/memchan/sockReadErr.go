package memchan

import "fmt"

// SockReadErr implements the SockReadErr interface
type SockReadErr struct {
	// closed is true when the error was caused by a graceful socket closure
	closed bool

	err error
}

// Error implements the Go error interface
func (err SockReadErr) Error() string {
	if err.closed {
		return "socket closed"
	}
	return fmt.Sprintf("reading socket failed: %s", err.err)
}

// IsCloseErr implements the SockReadErr interface
func (err SockReadErr) IsCloseErr() bool {
	return err.closed
}
