package test

import (
	"context"
	"testing"
	"time"

	"github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestServerMessageBufferOverflow tests sending messages to the server that are
// bigger than the servers message buffer
func TestServerMessageBufferOverflow(t *testing.T) {
	const messageBufferSize = uint32(2048)
	const messageHeaderSize = uint32(18) // type, identifier, name len, name

	// Determine the maximum payload length based on the message buffer size and
	// the size of the expected message header
	const maxPayloadSize = messageBufferSize - messageHeaderSize

	requestHandleCalled := tmdwg.NewTimedWaitGroup(1, 1*time.Second)
	correctRequestTriggeredHandler := tmdwg.NewTimedWaitGroup(1, 3*time.Second)

	// Generate a malicious payload that exceeds the message buffer size by 1
	// byte
	maliciousPayload := wwr.Payload{Data: make([]byte, maxPayloadSize+1)}

	// Initialize webwire server given only the request
	setup := setupTestServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				if string(msg.Name()) == "nooverfl" {
					correctRequestTriggeredHandler.Progress(1)
					return wwr.Payload{}, nil
				}

				requestHandleCalled.Progress(1)
				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{
			MessageBufferSize: messageBufferSize,
		},
	)

	// Initialize client
	client := setup.newClient(
		wwrclt.Options{
			// Disable autoconnect to avoid auto-reconnection on kick
			Autoconnect: wwr.Disabled,

			DefaultRequestTimeout: 2 * time.Second,
			// Use bigger buffers on the client
			MessageBufferSize: messageBufferSize * 2,
		},
		testClientHooks{},
	)

	// Wait until connected
	require.NoError(t, client.connection.Connect())

	// Send a perfectly sized request message and expect a reply
	replyCorrect, errCorrect := client.connection.Request(
		context.Background(),
		[]byte("nooverfl"),
		wwr.Payload{Data: make([]byte, maxPayloadSize)},
	)
	require.NotNil(t, replyCorrect)
	require.Nil(t, replyCorrect.Payload())
	require.NoError(t, errCorrect)
	replyCorrect.Close()

	require.NoError(t, correctRequestTriggeredHandler.Wait())

	// Send an overflowing request message and expect the server to close the
	// connection
	reply, err := client.connection.Request(
		context.Background(),
		[]byte("overflow"),
		maliciousPayload,
	)
	require.Nil(t, reply)
	require.Error(t, err)

	switch *argTransport {
	case "memchan":
		require.Equal(t, wwrclt.StatusConnected, client.connection.Status())
	default:
		require.Equal(t, wwrclt.StatusDisconnected, client.connection.Status())
	}

	// Expect the request handler not to be called
	require.Error(t, requestHandleCalled.Wait())
}
