package test

import (
	"context"
	"net/url"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	wwrfasthttp "github.com/qbeon/webwire-go/transport/fasthttp"
	"github.com/stretchr/testify/require"
)

// TestClientReqDisconnNoAutoconn tests disconnected error
// when trying to send a request while the server is unreachable
// and autoconn is disabled
func TestClientReqDisconnNoAutoconn(t *testing.T) {
	// Initialize client
	client, err := newClient(
		url.URL{Host: "127.0.0.1:65000"},
		&wwrfasthttp.Transport{},
		wwrclt.Options{
			Autoconnect:           wwr.Disabled,
			ReconnectionInterval:  5 * time.Millisecond,
			DefaultRequestTimeout: 50 * time.Millisecond,
		},
		testClientHooks{},
	)
	require.NoError(t, err)

	// Try to send a request and expect a DisconnectedErr error
	_, err = client.connection.Request(
		context.Background(),
		nil,
		wwr.Payload{Data: []byte("testdata")},
	)
	require.Error(t, err)
	require.IsType(t, wwr.DisconnectedErr{}, err)
	require.False(t, wwr.IsCanceledErr(err))
	require.False(t, wwr.IsTimeoutErr(err))
}
