package wwrerr

// ProtocolErr represents a protocol error
type ProtocolErr struct {
	Cause error
}

// NewProtocolErr constructs a new ProtocolErr error based on the actual error
func NewProtocolErr(err error) ProtocolErr {
	return ProtocolErr{
		Cause: err,
	}
}

// Error implements the error interface
func (err ProtocolErr) Error() string {
	return err.Cause.Error()
}
