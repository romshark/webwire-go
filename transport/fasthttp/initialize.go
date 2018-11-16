package fasthttp

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/fasthttp/websocket"
	wwr "github.com/qbeon/webwire-go"
	"github.com/valyala/fasthttp"
)

// Initialize implements the Transport interface
func (srv *Transport) Initialize(
	options wwr.ServerOptions,
	isShuttingdown wwr.IsShuttingDown,
	onNewConnection wwr.OnNewConnection,
) error {
	srv.isShuttingdown = isShuttingdown
	srv.onNewConnection = onNewConnection
	srv.readTimeout = options.ReadTimeout

	// Determine final address
	scheme := "http"
	if srv.TLS != nil {
		scheme = "https"
	}
	host := options.Host
	if host == "" {
		if srv.TLS != nil {
			host = ":https"
		} else {
			host = ":http"
		}
	}

	// Set default keep-alive period
	if srv.KeepAlive == 0 {
		srv.KeepAlive = 30 * time.Second
	}

	// Initialize TCP/IP listener
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return fmt.Errorf("TCP/IP listener setup failure: %s", err)
	}
	srv.addr = url.URL{
		Scheme: scheme,
		Host:   listener.Addr().String(),
		Path:   "/",
	}
	srv.listener = &tcpKeepAliveListener{
		listener.(*net.TCPListener),
		srv.KeepAlive,
	}

	// Set default HTTP server if none is specified
	if srv.HTTPServer == nil {
		srv.HTTPServer = &fasthttp.Server{
			Name: "webwire 2.0",
		}
	}

	srv.HTTPServer.ReadTimeout = options.ReadTimeout
	srv.HTTPServer.Handler = srv.handleAccept

	// Set default connection upgrader if none is specified
	if srv.Upgrader == nil {
		srv.Upgrader = &websocket.FastHTTPUpgrader{
			// Inherit buffer sizes from the HTTP server
			ReadBufferSize:  srv.HTTPServer.ReadBufferSize,
			WriteBufferSize: srv.HTTPServer.WriteBufferSize,
		}
	}

	// Create default loggers to std-out/err when no loggers are specified
	if srv.WarnLog == nil {
		srv.WarnLog = log.New(
			os.Stdout,
			"WWR_FASTHTTPWS_WARN: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		)
	}
	if srv.ErrorLog == nil {
		srv.ErrorLog = log.New(
			os.Stderr,
			"WWR_FASTHTTPWS_ERR: ",
			log.Ldate|log.Ltime|log.Lshortfile,
		)
	}

	return nil
}
