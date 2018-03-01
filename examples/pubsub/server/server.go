package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"unicode/utf16"

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

			// Encode message payload data to UTF16 for JavaScript compatibility
			runes := utf16.Encode([]rune(msg))
			data := make([]byte, len(runes)*2)
			dataItr := 0
			for i := 0; i < len(runes); i++ {
				tmp := make([]byte, 2)
				binary.LittleEndian.PutUint16(tmp, runes[i])
				data[dataItr] = tmp[0]
				dataItr++
				data[dataItr] = tmp[1]
				dataItr++
			}

			log.Printf("Broadcasting message '%s', to %d clients", msg, len(connectedClients))

			for client := range connectedClients {
				client.Signal("", webwire.Payload{
					Data: []byte(msg),
				})
			}
		}
	}()

	log.Printf("Listening on %s", server.Addr)
	server.Run()
}
