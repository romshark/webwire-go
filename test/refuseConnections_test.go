package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/qbeon/webwire-go/transport/memchan"
	"github.com/stretchr/testify/require"
)

// TestRefuseConnections tests refusing connections on the transport level
func TestRefuseConnections(t *testing.T) {
	numClients := 5

	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			Request: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				// Expect the following request to not even arrive
				t.Error("Not expected but reached")
				return wwr.Payload{}, nil
			},
		},
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

	clients := make([]*TestClient, numClients)
	for i := 0; i < numClients; i++ {
		clt := setup.NewClient(
			wwrclt.Options{
				DefaultRequestTimeout: 2 * time.Second,
				Autoconnect:           wwr.Disabled,
			},
			nil, // Use the default transport implementation
			TestClientHooks{},
		)
		clients[i] = clt
	}

	// Try sending requests
	for _, clt := range clients {
		_, err := clt.Connection.Request(
			context.Background(),
			[]byte("q"),
			wwr.Payload{},
		)
		require.Error(t, err)
		require.IsType(t, wwr.DisconnectedErr{}, err)
	}
}
