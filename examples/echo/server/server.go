package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	wwr "github.com/qbeon/webwire-go"
)

// EchoServer implements the webwire.ServerImplementation interface
type EchoServer struct{}

// OnOptions implements the webwire.ServerImplementation interface.
// Does nothing, not needed in this example
func (srv *EchoServer) OnOptions(_ http.ResponseWriter) {}

// OnSignal implements the webwire.ServerImplementation interface
// Does nothing, not needed in this example
func (srv *EchoServer) OnSignal(ctx context.Context) {}

// OnClientConnected implements the webwire.ServerImplementation interface.
// Does nothing, not needed in this example
func (srv *EchoServer) OnClientConnected(client *wwr.Client) {}

// OnClientDisconnected implements the webwire.ServerImplementation interface
// Does nothing, not needed in this example
func (srv *EchoServer) OnClientDisconnected(client *wwr.Client) {}

// BeforeUpgrade implements the webwire.ServerImplementation interface.
// Must return true to ensure incoming connections are accepted
func (srv *EchoServer) BeforeUpgrade(resp http.ResponseWriter, req *http.Request) bool {
	return true
}

// OnRequest implements the webwire.ServerImplementation interface.
// Returns the received message back to the client
func (srv *EchoServer) OnRequest(ctx context.Context) (response wwr.Payload, err error) {
	msg := ctx.Value(wwr.Msg).(wwr.Message)

	log.Printf("Replied to client: %s", msg.Client.RemoteAddr())

	// Reply to the request using the same data and encoding
	return msg.Payload, nil
}

var serverAddr = flag.String("addr", ":8081", "server address")

func main() {
	// Parse command line arguments
	flag.Parse()

	// Setup headed webwire server
	server, err := wwr.NewHeadedServer(
		&EchoServer{},
		wwr.HeadedServerOptions{
			ServerAddress: *serverAddr,
			ServerOptions: wwr.ServerOptions{
				WarnLog:  os.Stdout,
				ErrorLog: os.Stderr,
			},
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
	log.Printf("Listening on %s", server.Addr().String())
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}
}
