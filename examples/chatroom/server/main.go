package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	webwire "github.com/qbeon/webwire-go"
)

var serverAddr = flag.String("addr", ":8081", "server address")

func main() {
	// Parse command line arguments
	flag.Parse()

	// Initialize webwire server
	server, err := webwire.NewServer(
		*serverAddr,
		webwire.Hooks{
			OnClientConnected:    onClientConnected,
			OnClientDisconnected: onClientDisconnected,
			OnSignal:             onSignal,
			OnRequest:            onRequest,
			OnSessionCreated:     onSessionCreated,
			OnSessionLookup:      onSessionLookup,
			OnSessionClosed:      onSessionClosed,
		},
		os.Stdout, os.Stderr,
	)
	if err != nil {
		panic(fmt.Errorf(
			"Failed creating a new WebWire server instance: %s", err,
		))
	}

	// Launch server
	log.Printf("Listening on %s", server.Addr)
	server.Run()
}