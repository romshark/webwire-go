package gorilla

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/connopt"
)

// TLS represents TLS configuration
type TLS struct {
	CertFilePath       string
	PrivateKeyFilePath string
	Config             *tls.Config
}

// Transport implements the webwire transport layer with gorilla/websocket
type Transport struct {
	// OnOptions is invoked when the websocket endpoint is examined by the
	// client using the HTTP OPTION method.
	OnOptions func(resp http.ResponseWriter, req *http.Request)

	// BeforeUpgrade is invoked right before the upgrade of the connection of an
	// incoming HTTP request to a WebSocket connection and can be used to
	// intercept, configure or prevent incoming connections. BeforeUpgrade must
	// return the connection options to be applied or set options.Connection to
	// wwr.Refuse to refuse the incoming connection
	BeforeUpgrade func(
		resp http.ResponseWriter,
		req *http.Request,
	) connopt.ConnectionOptions

	// WarnLog defines the warn logging output target
	WarnLog *log.Logger

	// ErrorLog defines the error logging output target
	ErrorLog *log.Logger

	// KeepAlive enables the keep-alive option if set to a duration above -1.
	// KeepAlive is automatically set to 30 seconds when it's set to 0
	KeepAlive time.Duration

	// Upgrader specifies the websocket connection upgrader
	Upgrader *websocket.Upgrader

	// HTTPServer specifies the net/http server
	HTTPServer *http.Server

	// TLS enables TLS encryption if specified
	TLS *TLS

	listener        *tcpKeepAliveListener
	addr            url.URL
	readTimeout     time.Duration
	isShuttingdown  wwr.IsShuttingDown
	onNewConnection wwr.OnNewConnection
}

// Address returns the URL address the server is listening on
func (srv *Transport) Address() url.URL {
	return srv.addr
}
