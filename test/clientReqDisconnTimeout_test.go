package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientReqDisconnTimeout tests request timeout
// when the server is unreachable and autoconnect is enabled
func TestClientReqDisconnTimeout(t *testing.T) {
	// Initialize client
	client := newCallbackPoweredClient(
		"127.0.0.1:65000",
		wwrclt.Options{
			ReconnectionInterval:  5 * time.Millisecond,
			DefaultRequestTimeout: 50 * time.Millisecond,
		},
		callbackPoweredClientHooks{},
	)

	// Send request and await reply
	_, err := client.connection.Request(
		context.Background(),
		"",
		wwr.NewPayload(wwr.EncodingBinary, []byte("testdata")),
	)
	require.Error(t, err)
	require.IsType(t, wwr.TimeoutErr{}, err)
	require.True(t, wwr.IsTimeoutErr(err))
	require.False(t, wwr.IsCanceledErr(err))
}
