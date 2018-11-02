package test

import (
	"testing"
	"time"

	"github.com/qbeon/webwire-go/message"

	"github.com/stretchr/testify/require"

	"github.com/fasthttp/websocket"
	wwr "github.com/qbeon/webwire-go"
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

	readTimeout := 2 * time.Second

	// Setup a regular websocket connection
	serverAddr := server.AddressURL()
	if serverAddr.Scheme == "https" {
		serverAddr.Scheme = "wss"
	} else {
		serverAddr.Scheme = "ws"
	}

	conn, _, err := websocket.DefaultDialer.Dial(serverAddr.String(), nil)
	require.NoError(t, err)
	defer conn.Close()

	// Await the server configuration push message
	conn.SetReadDeadline(time.Now().Add(readTimeout))
	_, msgData, readErr := conn.ReadMessage()
	require.NoError(t, readErr)
	require.NotNil(t, msgData)

	// Parse server config message
	msg := &message.Message{}
	parsedMessageType, err := msg.Parse(msgData)
	require.True(t, parsedMessageType)
	require.NoError(t, err)
	require.Equal(t, [8]byte{0, 0, 0, 0, 0, 0, 0, 0}, msg.Identifier)
	require.Nil(t, msg.Name)
	require.Equal(t, message.ServerConfiguration{
		MajorProtocolVersion: 2,
		MinorProtocolVersion: 0,
		ReadTimeout:          serverReadTimeout,
		MessageBufferSize:    messageBufferSize,
	}, msg.ServerConfiguration)
}
