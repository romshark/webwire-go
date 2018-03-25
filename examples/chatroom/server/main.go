package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/examples/chatroom/server/state"
)

var serverAddr = flag.String("addr", ":9090", "server address")

func onClientConnected(newClient *wwr.Client) {
	state.State.AddConnected(newClient)
	log.Printf("New client connected: %s | %s", newClient.RemoteAddr(), newClient.UserAgent())
}

func onClientDisconnected(client *wwr.Client) {
	state.State.RemoveConnected(client)
	log.Printf("Client %s disconnected", client.RemoteAddr())
}

func main() {
	// Parse command line arguments
	flag.Parse()

	// Setup webwire server
	_, _, addr, runServer, stopServer, err := wwr.SetupServer(wwr.SetupOptions{
		ServerAddress: *serverAddr,
		ServerOptions: wwr.ServerOptions{
			SessionsEnabled: true,
			Hooks: wwr.Hooks{
				OnClientConnected:    onClientConnected,
				OnClientDisconnected: onClientDisconnected,
				OnRequest:            onRequest,
			},
			WarnLog:  os.Stdout,
			ErrorLog: os.Stderr,
		},
	})
	if err != nil {
		panic(fmt.Errorf("Failed setting up WebWire server: %s", err))
	}

	// Listen for OS signals and shutdown server in case of demanded termination
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-osSignals
		log.Printf("Termination demanded by the OS (%s), shutting down...", sig)
		if err := stopServer(); err != nil {
			log.Printf("Error during server shutdown: %s", err)
		}
		log.Println("Server gracefully terminated")
	}()

	// Launch server
	log.Printf("Listening on %s", addr)
	if err := runServer(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}
}
