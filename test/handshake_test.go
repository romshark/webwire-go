package test

import (
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	"github.com/stretchr/testify/require"
)

// TestHandshake tests the connection establishment handshake testing the server
// configuration push message
func TestHandshake(t *testing.T) {
	serverReadTimeout := 3 * time.Second
	messageBufferSize := uint32(1024 * 8)

	// Initialize webwire server
	setup := SetupTestServer(
		t,
		&ServerImpl{},
		wwr.ServerOptions{
			ReadTimeout:       serverReadTimeout,
			MessageBufferSize: messageBufferSize,
		},
		nil, // Use the default transport implementation
	)

	readTimeout := 5 * time.Second

	socket, err := setup.NewDisconnectedClientSocket()
	require.NoError(t, err)

	require.NoError(t, socket.Dial(time.Time{}))

	// Await the server configuration push message
	msg := message.NewMessage(messageBufferSize)
	require.NoError(t, socket.Read(msg, time.Now().Add(readTimeout)))

	require.Equal(t, [8]byte{}, msg.MsgIdentifier)
	require.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, msg.MsgIdentifierBytes)
	require.Nil(t, msg.MsgName)
	require.Equal(t, message.ServerConfiguration{
		MajorProtocolVersion: 2,
		MinorProtocolVersion: 0,
		ReadTimeout:          serverReadTimeout,
		MessageBufferSize:    messageBufferSize,
	}, msg.ServerConfiguration)
}
