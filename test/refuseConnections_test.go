package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	wwrfasthttp "github.com/qbeon/webwire-go/transport/fasthttp"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

// TestRefuseConnections tests refusal of connection before their upgrade to a
// websocket connection
func TestRefuseConnections(t *testing.T) {
	numClients := 5

	// Initialize server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				_ wwr.Connection,
				msg wwr.Message,
			) (wwr.Payload, error) {
				// Expect the following request to not even arrive
				t.Error("Not expected but reached")
				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{
			Transport: &wwrfasthttp.Transport{
				BeforeUpgrade: func(
					_ *fasthttp.RequestCtx,
				) wwr.ConnectionOptions {
					// Refuse all incoming connections
					return wwr.ConnectionOptions{Connection: wwr.Refuse}
				},
			},
		},
	)

	clients := make([]*callbackPoweredClient, numClients)
	for i := 0; i < numClients; i++ {
		clt := newCallbackPoweredClient(
			server.Address(),
			wwrclt.Options{
				DefaultRequestTimeout: 2 * time.Second,
				Autoconnect:           wwr.Disabled,
			},
			callbackPoweredClientHooks{},
		)
		defer clt.connection.Close()
		clients[i] = clt

		// Try connect
		require.Error(t, clt.connection.Connect())
	}

	// Try sending requests
	for i := 0; i < numClients; i++ {
		clt := clients[i]
		_, err := clt.connection.Request(
			context.Background(),
			[]byte("q"),
			wwr.Payload{},
		)
		require.Error(t, err)
		require.IsType(t, wwr.DisconnectedErr{}, err)
	}
}
