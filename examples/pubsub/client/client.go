package main

import (
	"flag"
	"log"
	"net/url"
	"sync"
	"time"

	wwr "github.com/qbeon/webwire-go"

	wwrclt "github.com/qbeon/webwire-go/client"
)

var serverAddr = flag.String("addr", ":8081", "server address")
var counterTarget = flag.Uint("n", 6, "number of signals to listen for")

// PubSubClient implements the wwrclt.Implementation interface
type PubSubClient struct {
	connection    wwrclt.Client
	target        uint
	counter       uint
	targetReached sync.WaitGroup
}

// NewPubSubClient constructs and returns a new pub-sub client instance
func NewPubSubClient(
	serverAddr url.URL,
	counterTarget uint,
) (*PubSubClient, error) {
	newPubSubClient := &PubSubClient{
		target:        counterTarget,
		counter:       0,
		targetReached: sync.WaitGroup{},
	}

	newPubSubClient.targetReached.Add(int(counterTarget))

	// Initialize connection
	connection, err := wwrclt.NewClient(
		serverAddr,
		newPubSubClient,
		wwrclt.Options{
			// Default timeout for timed requests
			DefaultRequestTimeout: 10 * time.Second,
			ReconnectionInterval:  2 * time.Second,
		},
		nil, // No TLS configuration
	)
	if err != nil {
		return nil, err
	}

	newPubSubClient.connection = connection

	return newPubSubClient, nil
}

// OnDisconnected implements the wwrclt.Implementation interface
func (clt *PubSubClient) OnDisconnected() {}

// OnSessionClosed implements the wwrclt.Implementation interface
func (clt *PubSubClient) OnSessionClosed() {}

// OnSessionCreated implements the wwrclt.Implementation interface
func (clt *PubSubClient) OnSessionCreated(_ *wwr.Session) {}

// OnSignal implements the wwrclt.Implementation interface
func (clt *PubSubClient) OnSignal(message wwr.Message) {
	clt.counter++
	log.Printf(
		"Signal %d of %d received: %s",
		clt.counter,
		clt.target,
		string(message.Payload().Data()),
	)
	clt.targetReached.Done()
}

// AwaitCounterTargetReached blocks the calling goroutine until the counter
// target is reached
func (clt *PubSubClient) AwaitCounterTargetReached() {
	clt.targetReached.Wait()
}

func main() {
	// Parse command line arguments
	flag.Parse()

	// Initialize a new pub-sub client instance
	client, err := NewPubSubClient(url.URL{Host: *serverAddr}, *counterTarget)
	if err != nil {
		panic(err)
	}

	// Wait until N signals are received before disconnecting
	client.AwaitCounterTargetReached()
}
