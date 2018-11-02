package test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientReqDisconnNoAutoconn tests disconnected error
// when trying to send a request while the server is unreachable
// and autoconn is disabled
func TestClientReqDisconnNoAutoconn(t *testing.T) {
	// Initialize client
	client := newCallbackPoweredClient(
		url.URL{Host: "127.0.0.1:65000"},
		wwrclt.Options{
			Autoconnect:           wwr.Disabled,
			ReconnectionInterval:  5 * time.Millisecond,
			DefaultRequestTimeout: 50 * time.Millisecond,
		},
		callbackPoweredClientHooks{},
	)

	// Try to send a request and expect a DisconnectedErr error
	_, err := client.connection.Request(
		context.Background(),
		nil,
		wwr.NewPayload(wwr.EncodingBinary, []byte("testdata")),
	)
	require.Error(t, err)
	require.IsType(t, wwr.DisconnectedErr{}, err)
	require.False(t, wwr.IsCanceledErr(err))
	require.False(t, wwr.IsTimeoutErr(err))
}
