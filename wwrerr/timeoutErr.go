package wwrerr

// TimeoutErr represents a failure caused by a timeout
type TimeoutErr struct {
	Cause error
}

// Error implements the error interface
func (err TimeoutErr) Error() string {
	return err.Cause.Error()
}
