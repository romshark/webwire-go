package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

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
		// Address of the webwire server
		*serverAddr,

		// No hooks required in this example
		webwireClient.Hooks{
			OnServerSignal: func(payload webwire.Payload) {
				counter++
				log.Printf("Signal %d received: %s", counter, string(payload.Data))
				todo.Done()
			},
		},

		// Default timeout for timed requests
		5*time.Second,

		// Log writers
		os.Stdout,
		os.Stderr,
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
