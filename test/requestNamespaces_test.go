package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestRequestNamespaces tests correct handling of namespaced requests
func TestRequestNamespaces(t *testing.T) {
	currentStep := 1

	shortestPossibleName := "s"
	buf := make([]rune, 255)
	for i := 0; i < 255; i++ {
		buf[i] = 'x'
	}
	longestPossibleName := string(buf)

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
					assert.Equal(t, "", msgName)
				case 2:
					assert.Equal(t, shortestPossibleName, msgName)
				case 3:
					assert.Equal(t, longestPossibleName, msgName)
				}
				return nil, nil
			},
		},
		webwire.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.AddressURL(),
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
		"",
		webwire.NewPayload(webwire.EncodingBinary, []byte("dummy")),
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
		webwire.NewPayload(webwire.EncodingBinary, []byte("dummy")),
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
		webwire.NewPayload(webwire.EncodingBinary, []byte("dummy")),
	)
	require.NoError(t, err)
}
