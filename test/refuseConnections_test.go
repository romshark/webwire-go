package test

import (
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/transport/memchan"
	"github.com/stretchr/testify/require"
)

// TestRefuseConnections tests refusing connections on the transport level
func TestRefuseConnections(t *testing.T) {
	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{},
		wwr.ServerOptions{},
		&memchan.Transport{
			OnBeforeCreation: func() wwr.ConnectionOptions {
				// Refuse all incoming connections
				return wwr.ConnectionOptions{
					Connection: wwr.Refuse,
				}
			},
		},
	)

	// Initialize client
	sock, err := setup.NewDisconnectedClientSocket()
	require.NoError(t, err)

	// Try connect
	require.Error(t, sock.Dial(time.Time{}))
}
