package test

import (
	"context"
	"testing"
	"time"

	"github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	"github.com/stretchr/testify/require"
)

// TestServerMessageBufferOverflow tests sending messages to the server that are
// bigger than the servers message buffer
func TestServerMessageBufferOverflow(t *testing.T) {
	const messageBufferSize = uint32(2048)
	const messageHeaderSize = uint32(10) // type, identifier, name len

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
		nil, // Use the default transport implementation
	)

	socket := setup.newClientSocket()

	// Await the server configuration push message
	confMsg := message.NewMessage(64)
	require.NoError(t, socket.Read(confMsg, time.Time{}))

	// Send an overflowing request message and expect the server to close the
	// connection
	writer, err := socket.GetWriter()
	require.NoError(t, err)
	require.NotNil(t, writer)

	require.NoError(t, message.WriteMsgRequest(
		writer,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0},
		[]byte{},
		maliciousPayload.Encoding,
		maliciousPayload.Data,
		true,
	))

	switch *argTransport {
	case "memchan":
	default:
		require.Error(t, socket.Read(confMsg, time.Time{}))
		require.False(t, socket.IsConnected())
	}

	// Expect the request handler not to be called
	require.Error(t, requestHandleCalled.Wait())
}
