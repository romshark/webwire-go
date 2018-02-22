package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	webwire "github.com/qbeon/webwire-go"
)

var serverAddr = flag.String("addr", ":8081", "server address")

var connectedClients = make(map[*webwire.Client]bool)

func onClientConnected(client *webwire.Client) {
	connectedClients[client] = true
}

func onClientDisconnected(client *webwire.Client) {
	delete(connectedClients, client)
}

func main() {
	// Parse command line arguments
	flag.Parse()

	// Initialize webwire server
	server, err := webwire.NewServer(
		*serverAddr,
		webwire.Hooks{
			OnClientConnected:    onClientConnected,
			OnClientDisconnected: onClientDisconnected,
		},
		os.Stdout, os.Stderr,
	)
	if err != nil {
		panic(fmt.Errorf(
			"Failed creating a new WebWire server instance: %s", err,
		))
	}

	// Begin sending the current time in 1 second intervals
	go func() {
		for {
			time.Sleep(1 * time.Second)

			if len(connectedClients) < 1 {
				log.Println("No clients connected, aborting broadcast")
				continue
			}

			msg := time.Now().String()

			log.Printf("Broadcasting message '%s', to %d clients", msg, len(connectedClients))

			for client, _ := range connectedClients {
				client.Signal([]byte(msg))
			}
		}
	}()

	log.Printf("Listening on %s", server.Addr)
	server.Run()
}
