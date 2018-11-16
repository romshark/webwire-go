package wwrerr

// SockReadErr defines the interface of a webwire.Socket.Read error
type SockReadErr interface {
	error

	// IsCloseErr must return true if the error represents a closure error
	IsCloseErr() bool
}
