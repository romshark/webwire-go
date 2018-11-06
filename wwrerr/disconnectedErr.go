package wwrerr

// DisconnectedErr represents an error indicating a disconnected connection
type DisconnectedErr struct {
	Cause error
}

// Error implements the error interface
func (err DisconnectedErr) Error() string {
	if err.Cause == nil {
		return "disconnected"
	}
	return err.Cause.Error()
}
