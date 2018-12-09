package test

import (
	"context"
	"sync"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/assert"
)

// TestNamedSignal tests correct handling of named signals
func TestNamedSignal(t *testing.T) {
	shortestNameSignalArrived := sync.WaitGroup{}
	shortestNameSignalArrived.Add(1)
	longestNameSignalArrived := sync.WaitGroup{}
	longestNameSignalArrived.Add(1)
	currentStep := 1

	shortestPossibleName := []byte("s")
	longestPossibleName := make([]byte, 255)
	for i := range longestPossibleName {
		longestPossibleName[i] = 'x'
	}

	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Signal: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) {
				msgName := msg.Name()
				switch currentStep {
				case 1:
					assert.Equal(t, shortestPossibleName, msgName)
					shortestNameSignalArrived.Done()
				case 2:
					assert.Equal(t, longestPossibleName, msgName)
					longestNameSignalArrived.Done()
				}
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	// Send request with the shortest possible name
	currentStep = 1
	signal(t, sock, shortestPossibleName, payload.Payload{})
	shortestNameSignalArrived.Wait()

	// Send request with the longest possible name
	currentStep = 2
	signal(t, sock, longestPossibleName, payload.Payload{})
	longestNameSignalArrived.Wait()
}
