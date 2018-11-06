package wwrerr

// SockReadErr defines the interface of a webwire.Socket.Read error
type SockReadErr interface {
	error

	// IsAbnormalCloseErr must return true if the error represents
	// an abnormal closure error
	IsAbnormalCloseErr() bool
}
