package test

import (
	"context"
	"testing"
	"time"

	"github.com/qbeon/webwire-go/message"

	wwr "github.com/qbeon/webwire-go"
	"github.com/stretchr/testify/require"
)

// TestRequestNoNameNoPayload tests sending requests without both a name and a
// payload expecting the server to reject the message
func TestRequestNoNameNoPayload(t *testing.T) {
	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				// Expect the following request to not even arrive
				t.Error("Not expected but reached")
				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	// TODO: improve test by avoiding the use of the client but performing the
	// request over a raw socket to ensure the client doesn't filter the request
	// out preemtively so it never even reaches the server

	// Send request without a name and without a payload.
	// Expect a protocol error in return not sending the invalid request off
	writer, err := sock.GetWriter()
	require.NoError(t, err)
	require.NotNil(t, writer)

	bytesWritten, err := writer.Write([]byte{
		message.MsgRequestBinary, // Type
		0, 0, 0, 0, 0, 0, 0, 0,   // Identifier
		0, // Name length
	})
	require.NoError(t, err)
	require.Equal(t, 10, bytesWritten)

	require.NoError(t, writer.Close())

	// Expect the socket to be closed by the server due to protocol violation
	msg := message.NewMessage(1024)
	readErr := sock.Read(msg, time.Time{})
	require.NotNil(t, readErr)
	require.True(t, readErr.IsCloseErr())
	require.False(t, sock.IsConnected())
}
