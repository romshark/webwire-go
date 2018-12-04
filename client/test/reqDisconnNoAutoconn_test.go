package client_test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	wwrtst "github.com/qbeon/webwire-go/test"
	"github.com/qbeon/webwire-go/transport/memchan"
	"github.com/stretchr/testify/require"
)

// TestReqDisconnNoAutoconn tests Client.Request expecting it to return a
// disconnected-error when trying to send a request while the server is
// unreachable and autoconn is disabled
func TestReqDisconnNoAutoconn(t *testing.T) {
	// Initialize client
	client, err := wwrtst.NewClient(
		wwrclt.Options{
			Autoconnect:           wwr.Disabled,
			ReconnectionInterval:  5 * time.Millisecond,
			DefaultRequestTimeout: 50 * time.Millisecond,
		},
		&memchan.ClientTransport{},
		wwrtst.TestClientHooks{},
	)
	require.NoError(t, err)

	// Try to send a request and expect a DisconnectedErr error
	_, err = client.Connection.Request(
		context.Background(),
		nil,
		wwr.Payload{Data: []byte("testdata")},
	)
	require.Error(t, err)
	require.IsType(t, wwr.DisconnectedErr{}, err)
	require.False(t, wwr.IsCanceledErr(err))
	require.False(t, wwr.IsTimeoutErr(err))
}
