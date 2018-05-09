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

func main() {
	// Parse command line arguments
	flag.Parse()

	// Setup headed webwire server
	server, err := wwr.NewHeadedServer(
		NewChatRoomServer(),
		wwr.HeadedServerOptions{
			ServerAddress: *serverAddr,
			ServerOptions: wwr.ServerOptions{
				SessionsEnabled: true,
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
	log.Printf("Listening on %s", server.Addr().String())
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}
}
