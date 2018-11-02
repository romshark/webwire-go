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

// TestClientRestSessDisconnTimeout tests autoconnect timeout
// when the server is unreachable
func TestClientRestSessDisconnTimeout(t *testing.T) {
	// Initialize client
	client := newCallbackPoweredClient(
		url.URL{Host: "127.0.0.1:65000"},
		wwrclt.Options{
			ReconnectionInterval:  5 * time.Millisecond,
			DefaultRequestTimeout: 50 * time.Millisecond,
		},
		callbackPoweredClientHooks{},
	)

	// Send request and await reply
	err := client.connection.RestoreSession(
		context.Background(),
		[]byte("inexistent_key"),
	)
	require.Error(t, err)
	require.IsType(t, wwr.TimeoutErr{}, err)
	require.True(t, wwr.IsTimeoutErr(err))
}
