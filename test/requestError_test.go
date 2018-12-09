package test

import (
	"context"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/require"
)

// TestRequestError tests server-side request errors properly failing the
// client-side requests
func TestRequestError(t *testing.T) {
	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
				_ context.Context,
				_ wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Fail the request by returning an error
				return wwr.Payload{Data: []byte("garbage")}, wwr.RequestErr{
					Code:    "SAMPLE_ERROR",
					Message: "Sample error message",
				}
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	rep := request(t, sock, 192, []byte("r"), payload.Payload{})
	require.Equal(t, message.MsgErrorReply, rep.MsgType)
	require.Equal(t, []byte("SAMPLE_ERROR"), rep.MsgName)
	require.Equal(t, payload.Binary, rep.MsgPayload.Encoding)
	require.Equal(t, []byte("Sample error message"), rep.MsgPayload.Data)
}
