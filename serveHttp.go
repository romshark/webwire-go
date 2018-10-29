package webwire

import (
	"net/http"
)

// ServeHTTP will make the server listen for incoming HTTP requests
// eventually trying to upgrade them to WebSocket connections
func (srv *server) ServeHTTP(
	resp http.ResponseWriter,
	req *http.Request,
) {
	// TODO: implement ServeHTTP support
	resp.WriteHeader(http.StatusNotImplemented)
}
