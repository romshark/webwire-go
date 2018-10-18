package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	wwr "github.com/qbeon/webwire-go"
)

var serverAddr = flag.String("addr", ":9090", "server address")
var certFile = flag.String(
	"sslcert",
	"./server.crt",
	"path to the SSL certificate file",
)
var privateKeyFile = flag.String(
	"sslkey",
	"./server.key",
	"path to the SSL private-key file",
)

func main() {
	// Parse command line arguments
	flag.Parse()

	// Setup a new webwire server instance
	server, err := wwr.NewServerSecure(
		NewChatRoomServer(),
		wwr.ServerOptions{
			Host: *serverAddr,
			WarnLog: log.New(
				os.Stdout,
				"WARN: ",
				log.Ldate|log.Ltime|log.Lshortfile,
			),
			ErrorLog: log.New(
				os.Stderr,
				"ERR: ",
				log.Ldate|log.Ltime|log.Lshortfile,
			),
		},
		*certFile,
		*privateKeyFile,
		nil,
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

	// Launch server
	log.Printf("Listening on %v", server.Address())
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}
}
