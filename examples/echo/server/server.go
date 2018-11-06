package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	wwr "github.com/qbeon/webwire-go"
	"github.com/valyala/fasthttp"
)

// EchoServer implements the webwire.ServerImplementation interface
type EchoServer struct{}

// OnOptions implements the webwire.ServerImplementation interface.
// Sets HTTP access control headers to satisfy CORS
func (srv *EchoServer) OnOptions(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "WEBWIRE")
}

// OnSignal implements the webwire.ServerImplementation interface
// Does nothing, not needed in this example
func (srv *EchoServer) OnSignal(
	_ context.Context,
	_ wwr.Connection,
	_ wwr.Message,
) {
}

// OnClientConnected implements the webwire.ServerImplementation interface.
// Does nothing, not needed in this example
func (srv *EchoServer) OnClientConnected(client wwr.Connection) {}

// OnClientDisconnected implements the webwire.ServerImplementation interface
// Does nothing, not needed in this example
func (srv *EchoServer) OnClientDisconnected(_ wwr.Connection, _ error) {}

// BeforeUpgrade implements the webwire.ServerImplementation interface
func (srv *EchoServer) BeforeUpgrade(
	_ *fasthttp.RequestCtx,
) wwr.ConnectionOptions {
	return wwr.ConnectionOptions{}
}

// OnRequest implements the webwire.ServerImplementation interface.
// Returns the received message back to the client
func (srv *EchoServer) OnRequest(
	_ context.Context,
	client wwr.Connection,
	message wwr.Message,
) (response wwr.Payload, err error) {
	log.Printf("Replied to client: %s", client.Info().RemoteAddr)

	// Reply to the request using the same data and encoding
	return wwr.Payload{
		Encoding: message.PayloadEncoding(),
		Data:     message.Payload(),
	}, nil
}

var serverAddr = flag.String("addr", ":8081", "server address")

func main() {
	// Parse command line arguments
	flag.Parse()

	// Setup a new webwire server instance
	server, err := wwr.NewServer(
		&EchoServer{},
		wwr.ServerOptions{
			Host: *serverAddr,
		},
	)
	if err != nil {
		panic(fmt.Errorf("Failed setting up WebWire server: %s", err))
	}

	// Listen for OS signals and shutdown server in case of demanded termination
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-osSignals
		log.Printf("Termination demanded by the OS (%s), shutting down...", sig)
		if err := server.Shutdown(); err != nil {
			log.Printf("Error during server shutdown: %s", err)
		}
		log.Println("Server gracefully terminated")
	}()

	// Launch echo server
	log.Printf("Listening on %s", server.Address())
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}
}
