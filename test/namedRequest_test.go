package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				_ webwire.Connection,
				msg webwire.Message,
			) (webwire.Payload, error) {
				msgName := msg.Name()
				switch currentStep {
				case 1:
					assert.Nil(t, msgName)
				case 2:
					assert.Equal(t, shortestPossibleName, msgName)
				case 3:
					assert.Equal(t, longestPossibleName, msgName)
				}
				return webwire.Payload{}, nil
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

	// Send unnamed request
	_, err := client.Connection.Request(
		context.Background(),
		nil,
		webwire.Payload{Data: []byte("dummy")},
	)
	require.NoError(t, err)

	// Send request with the shortest possible name
	currentStep = 2
	_, err = client.Connection.Request(
		context.Background(),
		shortestPossibleName,
		webwire.Payload{Data: []byte("dummy")},
	)
	require.NoError(t, err)

	// Send request with the longest possible name
	currentStep = 3
	_, err = client.Connection.Request(
		context.Background(),
		longestPossibleName,
		webwire.Payload{Data: []byte("dummy")},
	)
	require.NoError(t, err)
}
