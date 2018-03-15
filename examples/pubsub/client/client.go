package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/qbeon/webwire-go"

	webwireClient "github.com/qbeon/webwire-go/client"
)

var serverAddr = flag.String("addr", ":8081", "server address")
var number = flag.Uint("n", 6, "number of signals to listen for")

func main() {
	// Parse command line arguments
	flag.Parse()

	var todo sync.WaitGroup
	todo.Add(int(*number))
	counter := 0

	// Initialize client
	client := webwireClient.NewClient(
		*serverAddr,
		webwireClient.Options{
			// No hooks required in this example
			Hooks: webwireClient.Hooks{
				OnServerSignal: func(payload webwire.Payload) {
					counter++
					log.Printf("Signal %d received: %s", counter, string(payload.Data))
					todo.Done()
				},
			},
		},
	)

	// Close the client connection in the end
	defer client.Close()

	// Connect to the server
	if err := client.Connect(); err != nil {
		panic(fmt.Errorf("Couldn't connect to the server: %s", err))
	}
	log.Printf("Connected to %s", *serverAddr)

	// Wait until N signals are received before disconnecting
	todo.Wait()
}
