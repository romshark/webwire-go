package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"unicode/utf16"

	wwr "github.com/qbeon/webwire-go"
)

var serverAddr = flag.String("addr", ":8081", "server address")

var connectedClients = make(map[*wwr.Client]bool)
var mapLock = sync.Mutex{}

func onClientConnected(client *wwr.Client) {
	mapLock.Lock()
	connectedClients[client] = true
	mapLock.Unlock()
}

func onClientDisconnected(client *wwr.Client) {
	mapLock.Lock()
	delete(connectedClients, client)
	mapLock.Unlock()
}

func main() {
	// Parse command line arguments
	flag.Parse()

	// Setup headed webwire server
	server, err := wwr.NewHeadedServer(wwr.HeadedServerOptions{
		ServerAddress: *serverAddr,
		ServerOptions: wwr.ServerOptions{
			Hooks: wwr.Hooks{
				OnClientConnected:    onClientConnected,
				OnClientDisconnected: onClientDisconnected,
			},
			WarnLog:  os.Stdout,
			ErrorLog: os.Stderr,
		},
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

	// Launch server
	log.Printf("Listening on %s", server.Addr().String())
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}
}
