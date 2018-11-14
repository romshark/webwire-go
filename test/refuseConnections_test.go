package test

import (
	"context"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	wwrtrn "github.com/qbeon/webwire-go/transport"
	wwrfasthttp "github.com/qbeon/webwire-go/transport/fasthttp"
	wwrmemchan "github.com/qbeon/webwire-go/transport/memchan"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

// TestRefuseConnections tests refusal of connection before their upgrade to a
// websocket connection
func TestRefuseConnections(t *testing.T) {
	numClients := 5

	// Prepare transport layer implementation parameters
	var transImpl wwrtrn.Transport
	switch *argTransport {
	case "fasthttp/websocket":
		transImpl = &wwrfasthttp.Transport{
			BeforeUpgrade: func(
				_ *fasthttp.RequestCtx,
			) wwr.ConnectionOptions {
				// Refuse all incoming connections
				return wwr.ConnectionOptions{Connection: wwr.Refuse}
			},
		}
	case "memchan":
		transImpl = &wwrmemchan.Transport{
			// Refuse all incoming connections
			ConnectionOptions: wwr.ConnectionOptions{
				Connection: wwr.Refuse,
			},
		}
	default:
		t.Fatalf("unexpected transport implementation: %s", *argTransport)
	}

	// Initialize server
	setup := setupTestServer(
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
		wwr.ServerOptions{},
		transImpl,
	)

	clients := make([]*testClient, numClients)
	for i := 0; i < numClients; i++ {
		clt := setup.newClient(
			wwrclt.Options{
				DefaultRequestTimeout: 2 * time.Second,
				Autoconnect:           wwr.Disabled,
			},
			nil, // Use the default transport implementation
			testClientHooks{},
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
