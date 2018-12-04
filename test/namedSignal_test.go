package test

import (
	"context"
	"sync"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNamedSignal tests correct handling of named signals
func TestNamedSignal(t *testing.T) {
	unnamedSignalArrived := sync.WaitGroup{}
	unnamedSignalArrived.Add(1)
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
				_ webwire.Connection,
				msg webwire.Message,
			) {
				msgName := msg.Name()
				switch currentStep {
				case 1:
					assert.Nil(t, msgName)
					unnamedSignalArrived.Done()
				case 2:
					assert.Equal(t, shortestPossibleName, msgName)
					shortestNameSignalArrived.Done()
				case 3:
					assert.Equal(t, longestPossibleName, msgName)
					longestNameSignalArrived.Done()
				}
			},
		},
		webwire.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	client := setup.NewClient(
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		TestClientHooks{},
	)

	// Send unnamed signal
	require.NoError(t, client.Connection.Signal(
		context.Background(),
		nil, // No name
		webwire.Payload{Data: []byte("dummy")},
	))
	unnamedSignalArrived.Wait()

	// Send request with the shortest possible name
	currentStep = 2
	require.NoError(t, client.Connection.Signal(
		context.Background(),
		shortestPossibleName,
		webwire.Payload{Data: []byte("dummy")},
	))
	shortestNameSignalArrived.Wait()

	// Send request with the longest possible name
	currentStep = 3
	require.NoError(t, client.Connection.Signal(
		context.Background(),
		longestPossibleName,
		webwire.Payload{Data: []byte("dummy")},
	))
	longestNameSignalArrived.Wait()
}
