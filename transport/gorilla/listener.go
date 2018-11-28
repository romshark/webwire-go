package gorilla

import (
	"net"
	"time"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
	period time.Duration
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	if ln.period > 0 {
		if err := tc.SetKeepAlive(true); err != nil {
			return nil, err
		}
		if err := tc.SetKeepAlivePeriod(ln.period); err != nil {
			return nil, err
		}
	}
	return tc, nil
}
