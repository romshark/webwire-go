package wwrerr

// MaxSessConnsReachedErr represents an authentication error indicating that the
// given session already reached the maximum number of concurrent connections
type MaxSessConnsReachedErr struct{}

// Error implements the error interface
func (err MaxSessConnsReachedErr) Error() string {
	return "reached maximum number of concurrent session connections"
}
