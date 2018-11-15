package fasthttp

import (
	"crypto/tls"
	"fmt"
)

func (srv *Transport) serveTLS() error {
	var config *tls.Config
	if srv.TLS.Config != nil {
		config = srv.TLS.Config.Clone()
	} else {
		config = &tls.Config{}
	}

	if len(config.Certificates) < 1 {
		// Load and set TLS certificate if none is yet set
		cert, err := tls.LoadX509KeyPair(
			srv.TLS.CertFilePath,
			srv.TLS.PrivateKeyFilePath,
		)
		if err != nil {
			return fmt.Errorf("couldn't load TLS key-pair: %s", err)
		}

		config.Certificates = []tls.Certificate{cert}
	}

	// Launch HTTPS server
	if err := srv.HTTPServer.Serve(
		tls.NewListener(srv.listener, config),
	); err != nil {
		return fmt.Errorf("HTTPS server failure: %s", err)
	}

	return nil
}

// Serve implements the Transport interface
func (srv *Transport) Serve() error {
	if srv.TLS != nil {
		// Serve HTTPS
		return srv.serveTLS()
	}

	// Serve HTTP
	if err := srv.HTTPServer.Serve(srv.listener); err != nil {
		return fmt.Errorf("HTTP server failure: %s", err)
	}

	return nil
}
