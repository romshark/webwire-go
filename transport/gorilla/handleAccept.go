package gorilla

import (
	"net/http"

	"github.com/qbeon/webwire-go/connopt"
)

var methodNameOptions = []byte("OPTIONS")

func (srv *Transport) handleAccept(
	resp http.ResponseWriter,
	req *http.Request,
) {
	// Reject incoming connections during shutdown, pretend the server is
	// temporarily unavailable
	if srv.isShuttingdown() {
		http.Error(resp, "server shutting down", http.StatusServiceUnavailable)
		return
	}

	// Handle OPTION requests
	if req.Method == "OPTIONS" {
		if srv.OnOptions != nil {
			srv.OnOptions(resp, req)
		}
		return
	}

	connectionOptions := connopt.ConnectionOptions{
		Connection:       connopt.Accept,
		ConcurrencyLimit: 0,
	}
	if srv.BeforeUpgrade != nil {
		connectionOptions = srv.BeforeUpgrade(resp, req)
	}

	// Abort connection establishment if the connection was refused
	if connectionOptions.Connection != connopt.Accept {
		return
	}

	// Copy the user agent string
	conn, err := srv.Upgrader.Upgrade(resp, req, nil)
	if err != nil {
		// Establish connection
		srv.ErrorLog.Print("upgrade failed:", err)
		http.Error(
			resp,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	srv.handleConnection(
		connectionOptions,
		[]byte(req.UserAgent()),
		conn,
	)
}
