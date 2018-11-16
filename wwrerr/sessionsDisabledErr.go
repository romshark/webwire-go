package wwrerr

// SessionsDisabledErr represents an error indicating that the server has
// sessions disabled
type SessionsDisabledErr struct{}

// Error implements the error interface
func (err SessionsDisabledErr) Error() string {
	return "sessions are disabled for this server"
}
