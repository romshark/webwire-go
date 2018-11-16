package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrfasthttp "github.com/qbeon/webwire-go/transport/fasthttp"
)

var argServerAddr = flag.String("addr", ":9090", "server address")
var argCertFilePath = flag.String(
	"sslcert",
	"./server.crt",
	"path to the SSL certificate file",
)
var argPrivateKeyFile = flag.String(
	"sslkey",
	"./server.key",
	"path to the SSL private-key file",
)

func main() {
	// Parse command line arguments
	flag.Parse()

	// Setup a new webwire server instance
	server, err := wwr.NewServer(
		NewChatRoomServer(),
		wwr.ServerOptions{
			Host: *argServerAddr,
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
			ReadTimeout: 3 * time.Second,
		},
		&wwrfasthttp.Transport{
			TLS: &wwrfasthttp.TLS{
				CertFilePath:       *argCertFilePath,
				PrivateKeyFilePath: *argPrivateKeyFile,
				Config:             nil,
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

	// Launch server
	addr := server.Address()
	log.Printf("Listening on %s", addr.String())
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}
}
