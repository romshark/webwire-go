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

// TestRequestNamespaces tests correct handling of namespaced requests
func TestRequestNamespaces(t *testing.T) {
	currentStep := 1

	shortestPossibleName := []byte("s")
	longestPossibleName := make([]byte, 255)
	for i := range longestPossibleName {
		longestPossibleName[i] = 'x'
	}

	// Initialize server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
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
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Address(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	require.NoError(t, client.connection.Connect())

	/*****************************************************************\
		Step 1 - Unnamed request name
	\*****************************************************************/
	// Send unnamed request
	_, err := client.connection.Request(
		context.Background(),
		nil,
		webwire.Payload{Data: []byte("dummy")},
	)
	require.NoError(t, err)

	/*****************************************************************\
		Step 2 - Shortest possible request name
	\*****************************************************************/
	currentStep = 2
	// Send request with the shortest possible name
	_, err = client.connection.Request(
		context.Background(),
		shortestPossibleName,
		webwire.Payload{Data: []byte("dummy")},
	)
	require.NoError(t, err)

	/*****************************************************************\
		Step 3 - Longest possible request name
	\*****************************************************************/
	currentStep = 3
	// Send request with the longest possible name
	_, err = client.connection.Request(
		context.Background(),
		longestPossibleName,
		webwire.Payload{Data: []byte("dummy")},
	)
	require.NoError(t, err)
}
