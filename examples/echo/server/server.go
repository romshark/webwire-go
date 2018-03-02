package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	wwr "github.com/qbeon/webwire-go"
)

var serverAddr = flag.String("addr", ":8081", "server address")

func onRequest(ctx context.Context) (wwr.Payload, *wwr.Error) {
	msg := ctx.Value(wwr.Msg).(wwr.Message)
	client := msg.Client

	log.Printf("Replied to client: %s", client.RemoteAddr())

	// Reply to the request using the same data and encoding
	return msg.Payload, nil
}

func main() {
	// Parse command line arguments
	flag.Parse()

	// Setup webwire server
	_, _, addr, runServer, err := wwr.SetupServer(wwr.Options{
		Addr: *serverAddr,
		Hooks: wwr.Hooks{
			OnRequest: onRequest,
		},
		WarnLog:  os.Stdout,
		ErrorLog: os.Stderr,
	})
	if err != nil {
		panic(fmt.Errorf("Failed setting up WebWire server: %s", err))
	}

	log.Printf("Listening on %s", addr)

	if err := runServer(); err != nil {
		panic(fmt.Errorf("WebWire server failed: %s", err))
	}
}
