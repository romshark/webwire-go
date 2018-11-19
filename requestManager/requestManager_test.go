package requestmanager_test

import (
	"context"
	"errors"
	"testing"

	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/payload"
	reqman "github.com/qbeon/webwire-go/requestManager"
	"github.com/stretchr/testify/require"
)

// TestFulfillRequest tests RequestManager.Create, RequestManager.Fulfill,
// RequestManager.IsPending and Request.AwaitReply
func TestFulfillRequest(t *testing.T) {
	manager := reqman.NewRequestManager()

	// Create request
	request := manager.Create()
	require.NotNil(t, request)

	require.True(t, manager.IsPending(request.Identifier))

	// Fulfill the request
	pld := payload.Payload{
		Encoding: payload.Binary,
		Data:     []byte("test payload"),
	}
	require.True(t, manager.Fulfill(
		&message.Message{
			MsgIdentifier: request.Identifier,
			MsgPayload:    pld,
		},
	))

	require.False(t, manager.IsPending(request.Identifier))

	reply, err := request.AwaitReply(context.Background())
	require.NoError(t, err)
	require.NotNil(t, reply)
	require.Equal(t, pld.Encoding, reply.PayloadEncoding())
	require.Equal(t, pld.Data, reply.Payload())
}

// TestFailRequest tests RequestManager.Create, RequestManager.Fail,
// RequestManager.IsPending and Request.AwaitReply
func TestFailRequest(t *testing.T) {
	manager := reqman.NewRequestManager()

	// Create request
	request := manager.Create()
	require.NotNil(t, request)

	require.True(t, manager.IsPending(request.Identifier))

	// Fail the request
	manager.Fail(request.Identifier, errors.New("test error"))

	require.False(t, manager.IsPending(request.Identifier))

	reply, err := request.AwaitReply(context.Background())
	require.Nil(t, reply)
	require.Error(t, err)
}

// TestPendingRequests tests RequestManager.PendingRequests
func TestPendingRequests(t *testing.T) {
	manager := reqman.NewRequestManager()
	require.Equal(t, 0, manager.PendingRequests())

	// Create first request
	request1 := manager.Create()
	require.Equal(t, 1, manager.PendingRequests())

	// Create second request
	request2 := manager.Create()
	require.Equal(t, 2, manager.PendingRequests())

	// Fail the first request
	manager.Fail(request1.Identifier, errors.New("test error"))
	require.Equal(t, 1, manager.PendingRequests())

	// Fulfill the second request
	require.True(t, manager.Fulfill(
		&message.Message{
			MsgIdentifier: request2.Identifier,
			MsgPayload: payload.Payload{
				Encoding: payload.Binary,
				Data:     []byte("test payload"),
			},
		},
	))
	require.Equal(t, 0, manager.PendingRequests())
}
