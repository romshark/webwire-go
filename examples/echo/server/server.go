package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	webwire "github.com/qbeon/webwire-go"
)

var serverAddr = flag.String("addr", ":8081", "server address")

func onRequest(ctx context.Context) ([]byte, *webwire.Error) {
	msg := ctx.Value(webwire.MESSAGE).(webwire.Message)
	client := msg.Client

	log.Printf("Replied to client: %s", client.RemoteAddr())

	// Reply to the request
	return append([]byte("ECHO: "), msg.Payload...), nil
}

func main() {
	// Parse command line arguments
	flag.Parse()

	// Initialize webwire server
	server, err := webwire.NewServer(
		*serverAddr,
		webwire.Hooks{
			OnRequest: onRequest,
		},
		os.Stdout, os.Stderr,
	)
	if err != nil {
		panic(fmt.Errorf(
			"Failed creating a new WebWire server instance: %s", err,
		))
	}

	log.Printf("Listening on %s", server.Addr)

	server.Run()
}
