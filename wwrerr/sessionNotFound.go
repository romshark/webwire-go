package wwrerr

// SessionNotFoundErr represents a session restoration error indicating that the
// server didn't find the session to be restored
type SessionNotFoundErr struct{}

// Error implements the error interface
func (err SessionNotFoundErr) Error() string {
	return "session not found"
}
