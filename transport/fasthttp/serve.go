package fasthttp

import (
	"fmt"
	"net"
)

// Serve implements the Transport interface
func (srv *Transport) Serve() error {
	if srv.TLS != nil {
		if err := srv.HTTPServer.ServeTLS(
			tcpKeepAliveListener{srv.listener.(*net.TCPListener)},
			srv.TLS.CertFilePath,
			srv.TLS.PrivateKeyFilePath,
		); err != nil {
			return fmt.Errorf("HTTPS server failure: %s", err)
		}
	} else {
		if err := srv.HTTPServer.Serve(
			tcpKeepAliveListener{srv.listener.(*net.TCPListener)},
		); err != nil {
			return fmt.Errorf("HTTP server failure: %s", err)
		}
	}

	return nil
}
