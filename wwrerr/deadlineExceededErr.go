package wwrerr

// DeadlineExceededErr represents a failure due to an excess of a user-defined
// deadline
type DeadlineExceededErr struct {
	Cause error
}

// Error implements the error interface
func (err DeadlineExceededErr) Error() string {
	return err.Cause.Error()
}
