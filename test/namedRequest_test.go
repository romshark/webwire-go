package test

import (
	"context"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/assert"
)

// TestNamedRequest tests correct handling of named requests
func TestNamedRequest(t *testing.T) {
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
			Request: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				msgName := msg.Name()
				switch currentStep {
				case 1:
					assert.Equal(t, shortestPossibleName, msgName)
				case 2:
					assert.Equal(t, longestPossibleName, msgName)
				}
				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	// Send request with the shortest possible name
	currentStep = 1
	requestSuccess(t, sock, 32, shortestPossibleName, payload.Payload{})

	// Send request with the longest possible name
	currentStep = 2
	requestSuccess(t, sock, 32, longestPossibleName, payload.Payload{})
}
