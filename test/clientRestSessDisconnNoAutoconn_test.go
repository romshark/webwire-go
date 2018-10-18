package test

import (
	"net/url"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/require"
)

// TestClientRestSessDisconnNoAutoconn tests disconnected error
// when trying to manually restore the session
// while the server is unreachable and autoconn is disabled
func TestClientRestSessDisconnNoAutoconn(t *testing.T) {
	// Initialize client
	client := newCallbackPoweredClient(
		url.URL{Host: "127.0.0.1:65000"},
		wwrclt.Options{
			Autoconnect:           wwr.Disabled,
			ReconnectionInterval:  5 * time.Millisecond,
			DefaultRequestTimeout: 50 * time.Millisecond,
		},
		callbackPoweredClientHooks{},
		nil, // No TLS configuration
	)

	// Try to restore a session and expect a DisconnectedErr error
	err := client.connection.RestoreSession([]byte("inexistent_key"))
	require.Error(t, err)
	require.IsType(t, wwr.DisconnectedErr{}, err)
}
