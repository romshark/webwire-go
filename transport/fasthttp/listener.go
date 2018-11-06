package fasthttp

import (
	"net"
	"time"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept accepts incoming client connections
func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
