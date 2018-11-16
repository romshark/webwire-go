package wwrerr

// DialTimeoutErr represents a dialing error caused by a timeout
type DialTimeoutErr struct{}

// Error implements the error interface
func (err DialTimeoutErr) Error() string {
	return "dial timed out"
}
