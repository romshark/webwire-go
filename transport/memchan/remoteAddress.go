package memchan

import (
	"fmt"
)

// RemoteAddress represents a net.Addr interface implementation
type RemoteAddress struct {
	serverSocket *Socket
}

// Network implements the net.Addr interface
func (addr RemoteAddress) Network() string {
	return "memchan"
}

// String implements the net.Addr interface
func (addr RemoteAddress) String() string {
	return fmt.Sprintf("%p", addr.serverSocket)
}
