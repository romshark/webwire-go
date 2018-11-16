package wwrerr

// RequestErr represents an error returned when a request couldn't be processed
type RequestErr struct {
	Code    string
	Message string
}

// Error implements the error interface
func (err RequestErr) Error() string {
	return err.Message
}
