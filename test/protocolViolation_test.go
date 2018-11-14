package test

import (
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	"github.com/stretchr/testify/require"
)

// TestProtocolViolation tests sending messages that violate the protocol
func TestProtocolViolation(t *testing.T) {
	// Initialize webwire server
	setup := setupTestServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	defaultReadTimeout := 2 * time.Second

	// Setup a regular websocket connection
	try := func(m []byte) {
		socket := setup.newClientSocket()

		// Ignore the server configuration push-message
		confMsg := message.NewMessage(256)
		require.NoError(t, socket.Read(
			confMsg,
			time.Now().Add(defaultReadTimeout),
		))

		// Get writer
		writer, err := socket.GetWriter()
		require.NoError(t, err)

		// Write the message
		bytesWritten, writeErr := writer.Write(m)
		require.NoError(t, writeErr)
		require.Equal(t, len(m), bytesWritten)
		require.NoError(t, writer.Close())

		emptyMsg := message.NewMessage(256)
		readErr := socket.Read(emptyMsg, time.Now().Add(defaultReadTimeout))
		require.Error(t, readErr)

		require.False(t, socket.IsConnected())
	}

	// Test a message with an invalid type identifier (200, which is undefined)
	// and expect the server to ignore it returning no answer
	try([]byte{byte(200)})

	// Test a message with an invalid name length flag (bigger than name)
	// and expect the server to return a protocol violation error response
	try([]byte{
		message.MsgRequestBinary, // Message type identifier
		0, 0, 0, 0, 0, 0, 0, 0,   // Request identifier
		3,     // Name length flag
		0x041, // Name
	})
}
