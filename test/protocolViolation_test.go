package test

import (
	"testing"
	"time"

	"github.com/qbeon/webwire-go/message"

	"github.com/stretchr/testify/require"

	"github.com/fasthttp/websocket"
	wwr "github.com/qbeon/webwire-go"
)

// TestProtocolViolation tests sending messages that violate the protocol
func TestProtocolViolation(t *testing.T) {
	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{},
	)

	defaultReadTimeout := 2 * time.Second

	// Setup a regular websocket connection
	setupAndSend := func(
		message []byte,
	) (response []byte, writeErr, readErr error) {
		serverAddr := server.AddressURL()
		if serverAddr.Scheme == "https" {
			serverAddr.Scheme = "wss"
		} else {
			serverAddr.Scheme = "ws"
		}

		conn, _, err := websocket.DefaultDialer.Dial(serverAddr.String(), nil)
		require.NoError(t, err)
		defer conn.Close()

		// Ignore the server configuration push-message
		conn.SetReadDeadline(time.Now().Add(defaultReadTimeout))
		_, _, err = conn.ReadMessage()
		require.NoError(t, err)

		// Write the message
		writeErr = conn.WriteMessage(websocket.BinaryMessage, message)
		if writeErr != nil {
			return nil, writeErr, nil
		}

		// Await the response
		conn.SetReadDeadline(time.Now().Add(defaultReadTimeout))
		_, response, readErr = conn.ReadMessage()
		if readErr != nil {
			return nil, nil, readErr
		}

		return response, nil, nil
	}

	// Test a message with an invalid type identifier (200, which is undefined)
	// and expect the server to ignore it returning no answer
	func() {
		msg := []byte{byte(200)}
		response, writeErr, readErr := setupAndSend(msg)
		require.NoError(t, writeErr)
		require.Error(t, readErr)
		require.Nil(t, response)
	}()

	// Test a message with an invalid name length flag (bigger than name)
	// and expect the server to return a protocol violation error response
	func() {
		msg := []byte{
			message.MsgRequestBinary, // Message type identifier
			0, 0, 0, 0, 0, 0, 0, 0,   // Request identifier
			3,     // Name length flag
			0x041, // Name
		}
		response, writeErr, readErr := setupAndSend(msg)
		require.NoError(t, writeErr)
		require.NoError(t, readErr)
		require.Equal(t, []byte{
			message.MsgReplyProtocolError, // Message type identifier
			0, 0, 0, 0, 0, 0, 0, 0,        // Request identifier
		}, response)
	}()
}
