package fasthttp

import (
	"crypto/tls"
	"log"
	"net"
	"net/url"
	"time"

	"github.com/fasthttp/websocket"
	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/connopt"
	"github.com/valyala/fasthttp"
)

// TLS represents TLS configuration
type TLS struct {
	CertFilePath       string
	PrivateKeyFilePath string
	Config             *tls.Config
}

// Transport implements the webwire transport layer with fasthttp
type Transport struct {
	// OnOptions is invoked when the websocket endpoint is examined by the
	// client using the HTTP OPTION method.
	OnOptions func(*fasthttp.RequestCtx)

	// BeforeUpgrade is invoked right before the upgrade of the connection of an
	// incoming HTTP request to a WebSocket connection and can be used to
	// intercept, configure or prevent incoming connections. BeforeUpgrade must
	// return the connection options to be applied or set options.Connection to
	// wwr.Refuse to refuse the incoming connection
	BeforeUpgrade func(ctx *fasthttp.RequestCtx) connopt.ConnectionOptions

	// WarnLog defines the warn logging output target
	WarnLog *log.Logger

	// ErrorLog defines the error logging output target
	ErrorLog *log.Logger

	// Upgrader specifies the websocket connection upgrader
	Upgrader *websocket.FastHTTPUpgrader

	// HTTPServer specifies the FastHTTP server
	HTTPServer *fasthttp.Server

	// TLS enables TLS encryption if specified
	TLS *TLS

	listener        net.Listener
	addr            url.URL
	readTimeout     time.Duration
	isShuttingdown  wwr.IsShuttingDown
	onNewConnection wwr.OnNewConnection
}

// Address returns the URL address the server is listening on
func (srv *Transport) Address() url.URL {
	return srv.addr
}
