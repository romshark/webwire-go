package test

import (
	"context"
	"fmt"
	"testing"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/require"
)

// TestRequestInternalError tests returning non-ReqErr errors from the request
// handler
func TestRequestInternalError(t *testing.T) {
	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
				_ context.Context,
				_ wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Fail the request by returning a non-ReqErr error
				return wwr.Payload{Data: []byte("garbage")}, fmt.Errorf(
					"don't worry, this internal error is expected",
				)
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	sock, _ := setup.NewClientSocket()

	rep := request(t, sock, 192, []byte("r"), payload.Payload{})
	require.Equal(t, message.MsgInternalError, rep.MsgType)
	require.Nil(t, rep.MsgName)
	require.Equal(t, payload.Binary, rep.MsgPayload.Encoding)
	require.Nil(t, rep.MsgPayload.Data)
}
