package client

// DialTimeout represents a dialing error caused by a timeout
type DialTimeout struct{}

// Error implements the standard error interface
func (err DialTimeout) Error() string {
	return "dialing timed out"
}
