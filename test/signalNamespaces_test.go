package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	tmdwg "github.com/qbeon/tmdwg-go"
	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestSignalNamespaces tests correct handling of namespaced signals
func TestSignalNamespaces(t *testing.T) {
	unnamedSignalArrived := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	shortestNameSignalArrived := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	longestNameSignalArrived := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
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
			onSignal: func(
				_ context.Context,
				_ webwire.Connection,
				msg webwire.Message,
			) {
				msgName := msg.Name()
				switch currentStep {
				case 1:
					assert.Nil(t, msgName)
					unnamedSignalArrived.Progress(1)
				case 2:
					assert.Equal(t, shortestPossibleName, msgName)
					shortestNameSignalArrived.Progress(1)
				case 3:
					assert.Equal(t, longestPossibleName, msgName)
					longestNameSignalArrived.Progress(1)
				}
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
		Step 1 - Unnamed signal name
	\*****************************************************************/
	// Send unnamed signal
	require.NoError(t, client.connection.Signal(
		context.Background(),
		nil, // No name
		webwire.NewPayload(
			webwire.EncodingBinary,
			[]byte("dummy"),
		),
	))
	require.NoError(t,
		unnamedSignalArrived.Wait(),
		"Unnamed signal didn't arrive",
	)

	/*****************************************************************\
		Step 2 - Shortest possible request name
	\*****************************************************************/
	currentStep = 2

	// Send request with the shortest possible name
	require.NoError(t, client.connection.Signal(
		context.Background(),
		shortestPossibleName,
		webwire.NewPayload(webwire.EncodingBinary, []byte("dummy")),
	))
	require.NoError(t,
		shortestNameSignalArrived.Wait(),
		"Signal with shortest name didn't arrive",
	)

	/*****************************************************************\
		Step 3 - Longest possible request name
	\*****************************************************************/
	currentStep = 3

	// Send request with the longest possible name
	require.NoError(t, client.connection.Signal(
		context.Background(),
		longestPossibleName,
		webwire.NewPayload(webwire.EncodingBinary, []byte("dummy")),
	))
	require.NoError(t,
		longestNameSignalArrived.Wait(),
		"Signal with longest name didn't arrive",
	)
}
