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

// TestClientReqDisconnTimeout tests request timeout
// when the server is unreachable and autoconnect is enabled
func TestClientReqDisconnTimeout(t *testing.T) {
	// Initialize client
	client, err := newClient(
		url.URL{Host: "127.0.0.1:65000"},
		&wwrfasthttp.Transport{},
		wwrclt.Options{
			ReconnectionInterval:  5 * time.Millisecond,
			DefaultRequestTimeout: 50 * time.Millisecond,
		},
		testClientHooks{},
	)
	require.NoError(t, err)

	// Send request and await reply
	reply, err := client.connection.Request(
		context.Background(),
		nil,
		wwr.Payload{Data: []byte("testdata")},
	)
	require.Error(t, err)
	require.Nil(t, reply)
	require.IsType(t, wwr.TimeoutErr{}, err)
	require.True(t, wwr.IsTimeoutErr(err))
	require.False(t, wwr.IsCanceledErr(err))
}
