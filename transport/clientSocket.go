package transport

import (
	"net/url"
	"time"
)

// ClientSocket defines the abstract client socket implementation interface
type ClientSocket interface {
	// Dial connects the socket to the specified server
	Dial(serverAddr url.URL, deadline time.Time) error

	Socket
}
