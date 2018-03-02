package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"unicode/utf16"

	wwr "github.com/qbeon/webwire-go"
)

var serverAddr = flag.String("addr", ":8081", "server address")

var connectedClients = make(map[*wwr.Client]bool)

func onClientConnected(client *wwr.Client) {
	connectedClients[client] = true
}

func onClientDisconnected(client *wwr.Client) {
	delete(connectedClients, client)
}

func main() {
	// Parse command line arguments
	flag.Parse()

	// Setup webwire server
	_, _, addr, runServer, err := wwr.SetupServer(wwr.Options{
		Addr: *serverAddr,
		Hooks: wwr.Hooks{
			OnClientConnected:    onClientConnected,
			OnClientDisconnected: onClientDisconnected,
		},
		WarnLog:  os.Stdout,
		ErrorLog: os.Stderr,
	})
	if err != nil {
		panic(fmt.Errorf("Failed setting up WebWire server: %s", err))
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
				client.Signal("", wwr.Payload{
					Data: []byte(msg),
				})
			}
		}
	}()

	log.Printf("Listening on %s", addr)

	if err := runServer(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}
}
