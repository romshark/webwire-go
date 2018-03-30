package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	wwr "github.com/qbeon/webwire-go"
)

var serverAddr = flag.String("addr", ":8081", "server address")

func onRequest(ctx context.Context) (wwr.Payload, error) {
	msg := ctx.Value(wwr.Msg).(wwr.Message)
	client := msg.Client

	log.Printf("Replied to client: %s", client.RemoteAddr())

	// Reply to the request using the same data and encoding
	return msg.Payload, nil
}

func main() {
	// Parse command line arguments
	flag.Parse()

	// Setup headed webwire server
	server, err := wwr.NewHeadedServer(wwr.HeadedServerOptions{
		ServerAddress: *serverAddr,
		ServerOptions: wwr.ServerOptions{
			Hooks: wwr.Hooks{
				OnRequest: onRequest,
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
		if err := server.Shutdown(); err != nil {
			log.Printf("Error during server shutdown: %s", err)
		}
		log.Println("Server gracefully terminated")
	}()

	// Launch echo server
	log.Printf("Listening on %s", server.Addr().String())
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}
}
