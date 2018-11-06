package wwrerr

// CanceledErr represents a failure due to cancellation
type CanceledErr struct {
	Cause error
}

// Error implements the error interface
func (err CanceledErr) Error() string {
	return err.Cause.Error()
}

// IsTimeoutErr returns true if the given error is either a TimeoutErr
// or a DeadlineExceededErr, otherwise returns false
func IsTimeoutErr(err error) bool {
	switch err.(type) {
	case TimeoutErr:
		return true
	case DeadlineExceededErr:
		return true
	}
	return false
}

// IsCanceledErr returns true if the given error is a CanceledErr,
// otherwise returns false
func IsCanceledErr(err error) bool {
	switch err.(type) {
	case CanceledErr:
		return true
	}
	return false
}
