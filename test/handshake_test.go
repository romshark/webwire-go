package test

import (
	"testing"
	"time"

	"github.com/fasthttp/websocket"
	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	"github.com/stretchr/testify/require"
)

// TestHandshake tests the connection establishment handshake testing the server
// configuration push message
func TestHandshake(t *testing.T) {
	serverReadTimeout := 3 * time.Second
	messageBufferSize := uint32(1024 * 64)

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{
			ReadTimeout:       serverReadTimeout,
			MessageBufferSize: messageBufferSize,
		},
	)

	readTimeout := 5 * time.Second

	// Setup a regular websocket connection
	serverAddr := server.Address()
	if serverAddr.Scheme == "https" {
		serverAddr.Scheme = "wss"
	} else {
		serverAddr.Scheme = "ws"
	}

	conn, _, err := websocket.DefaultDialer.Dial(serverAddr.String(), nil)
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer conn.Close()

	// Await the server configuration push message
	require.NoError(t, conn.SetReadDeadline(time.Now().Add(readTimeout)))

	messageType, reader, err := conn.NextReader()
	require.NoError(t, err)
	require.Equal(t, websocket.BinaryMessage, messageType)

	msg := message.NewMessage(messageBufferSize)
	parsedMessageType, readErr := msg.Read(reader)
	require.NoError(t, readErr)
	require.True(t, parsedMessageType)

	require.Equal(t, [8]byte{0, 0, 0, 0, 0, 0, 0, 0}, msg.MsgIdentifier)
	require.Nil(t, msg.MsgName)
	require.Equal(t, message.ServerConfiguration{
		MajorProtocolVersion: 2,
		MinorProtocolVersion: 0,
		ReadTimeout:          serverReadTimeout,
		MessageBufferSize:    messageBufferSize,
	}, msg.ServerConfiguration)
}
