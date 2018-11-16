package wwrerr

// BufferOverflowErr represents a message buffer overflow error
type BufferOverflowErr struct{}

// Error implements the error interface
func (err BufferOverflowErr) Error() string {
	return "message buffer overflow"
}
