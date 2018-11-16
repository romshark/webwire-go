package wwrerr

// InternalErr represents a server-side internal error
type InternalErr struct{}

// Error implements the error interface
func (err InternalErr) Error() string {
	return "internal server error"
}
