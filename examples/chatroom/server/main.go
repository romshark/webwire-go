package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	wwr "github.com/qbeon/webwire-go"
)

var serverAddr = flag.String("addr", ":8081", "server address")

func main() {
	// Parse command line arguments
	flag.Parse()

	// Setup webwire server
	_, _, addr, runServer, err := wwr.SetupServer(wwr.Options{
		Addr: *serverAddr,
		Hooks: wwr.Hooks{
			OnClientConnected:    onClientConnected,
			OnClientDisconnected: onClientDisconnected,
			OnSignal:             onSignal,
			OnRequest:            onRequest,
			OnSessionCreated:     onSessionCreated,
			OnSessionLookup:      onSessionLookup,
			OnSessionClosed:      onSessionClosed,
		},
		WarnLog:  os.Stdout,
		ErrorLog: os.Stderr,
	})
	if err != nil {
		panic(fmt.Errorf("Failed setting up WebWire server: %s", err))
	}

	// Launch server
	log.Printf("Listening on %s", addr)
	if err := runServer(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}
}
