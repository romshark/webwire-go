package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	wwr "github.com/qbeon/webwire-go"
)

// PubSubServer implements the webwire.ServerImplementation interface
type PubSubServer struct {
	broadcastInterval time.Duration
	connectedClients  map[*wwr.Client]bool
	mapLock           sync.Mutex
}

// NewPubSubServer constructs a new pub-sub webwire server implementation instance
func NewPubSubServer() *PubSubServer {
	return &PubSubServer{
		1 * time.Second,
		make(map[*wwr.Client]bool),
		sync.Mutex{},
	}
}

// OnOptions implements the webwire.ServerImplementation interface.
// Sets HTTP access control headers to satisfy CORS
func (srv *PubSubServer) OnOptions(resp http.ResponseWriter) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "WEBWIRE")
}

// OnSignal implements the webwire.ServerImplementation interface
// Does nothing, not needed in this example
func (srv *PubSubServer) OnSignal(
	_ context.Context,
	_ *wwr.Client,
	_ *wwr.Message,
) {
}

// OnRequest implements the webwire.ServerImplementation interface.
// Does nothing, not needed in this example
func (srv *PubSubServer) OnRequest(
	_ context.Context,
	_ *wwr.Client,
	_ *wwr.Message,
) (response wwr.Payload, err error) {
	return wwr.Payload{}, wwr.ReqErr{
		Code:    "REQ_NOT_SUPPORTED",
		Message: "Requests are not supported on this server",
	}
}

// BeforeUpgrade implements the webwire.ServerImplementation interface.
// Must return true to ensure incoming connections are accepted
func (srv *PubSubServer) BeforeUpgrade(resp http.ResponseWriter, req *http.Request) bool {
	return true
}

// OnClientConnected implements the webwire.ServerImplementation interface.
// Registers a new connected client
func (srv *PubSubServer) OnClientConnected(client *wwr.Client) {
	srv.mapLock.Lock()
	srv.connectedClients[client] = true
	srv.mapLock.Unlock()
}

// OnClientDisconnected implements the webwire.ServerImplementation interface
// Deregisters a gone client
func (srv *PubSubServer) OnClientDisconnected(client *wwr.Client) {
	srv.mapLock.Lock()
	delete(srv.connectedClients, client)
	srv.mapLock.Unlock()
}

// Broadcast begins sending the current time in 1 second intervals.
// Blocks the calling goroutine
func (srv *PubSubServer) Broadcast() {
	for {
		time.Sleep(1 * time.Second)

		srv.mapLock.Lock()

		if len(srv.connectedClients) < 1 {
			log.Println("No clients connected, aborting broadcast")
			srv.mapLock.Unlock()
			continue
		}

		msg := time.Now().String()

		log.Printf("Broadcasting message '%s', to %d clients", msg, len(srv.connectedClients))

		for client := range srv.connectedClients {
			client.Signal("", wwr.Payload{
				Data: []byte(msg),
			})
		}
		srv.mapLock.Unlock()
	}
}

// Accept -addr CLI parameter defining the server address, default to :8081
var serverAddr = flag.String("addr", ":8081", "server address")

func main() {
	// Parse command line arguments
	flag.Parse()

	// Create a new webwire server implementation instance
	serverImpl := NewPubSubServer()

	// Setup headed webwire server
	server, err := wwr.NewHeadedServer(
		serverImpl,
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

	// Start broadcast
	go serverImpl.Broadcast()

	// Listen for OS signals and shutdown server in case of demanded termination
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		// Await OS signal
		sig := <-osSignals

		log.Printf("Termination demanded by the OS (%s), shutting down...", sig)

		// Shutdown the webwire server
		if err := server.Shutdown(); err != nil {
			log.Printf("Error during server shutdown: %s", err)
		}
		log.Println("Server gracefully terminated")
	}()

	log.Printf("Listening on %s", server.Addr().String())

	// Launch server
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}
}
