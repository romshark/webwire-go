package gorilla

import (
	"context"
	"fmt"
)

// Shutdown implements the Transport interface
func (srv *Transport) Shutdown() error {
	if srv.HTTPServer == nil {
		return nil
	}
	if err := srv.HTTPServer.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("couldn't properly shutdown the HTTP server: %s", err)
	}
	return nil
}
