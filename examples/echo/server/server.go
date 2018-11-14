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
	wwrfasthttp "github.com/qbeon/webwire-go/transport/fasthttp"
)

// EchoServer implements the webwire.ServerImplementation interface
type EchoServer struct{}

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
		&wwrfasthttp.Transport{},
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
	addr := server.Address()
	log.Printf("Listening on %s", addr.String())
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}
}
