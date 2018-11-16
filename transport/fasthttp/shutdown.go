package fasthttp

import "fmt"

// Shutdown implements the Transport interface
func (srv *Transport) Shutdown() error {
	if srv.HTTPServer == nil {
		return nil
	}
	if err := srv.HTTPServer.Shutdown(); err != nil {
		return fmt.Errorf("couldn't properly shutdown the HTTP server: %s", err)
	}
	return nil
}
