package wwrerr

// ServerShutdownErr represents a request error indicating that the request
// cannot be processed due to the server currently being shut down
type ServerShutdownErr struct{}

// Error implements the error interface
func (err ServerShutdownErr) Error() string {
	return "server is currently being shut down and won't process the request"
}
