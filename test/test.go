package test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	wwr "github.com/qbeon/webwire-go"
)

// setupServer helps setting up and launching the server
// together with the hosting http server
// setting up a headed server on a randomly assigned port
func setupServer(
	t *testing.T,
	impl *serverImpl,
	opts wwr.ServerOptions,
) wwr.Server {
	// Setup headed server on arbitrary port

	if impl.beforeUpgrade == nil {
		impl.beforeUpgrade = func(_ http.ResponseWriter, _ *http.Request) bool {
			return true
		}
	}
	if impl.onClientConnected == nil {
		impl.onClientConnected = func(_ wwr.Connection) {}
	}
	if impl.onClientDisconnected == nil {
		impl.onClientDisconnected = func(_ wwr.Connection) {}
	}
	if impl.onSignal == nil {
		impl.onSignal = func(
			_ context.Context,
			_ wwr.Connection,
			_ wwr.Message,
		) {
		}
	}
	if impl.onRequest == nil {
		impl.onRequest = func(
			_ context.Context,
			_ wwr.Connection,
			_ wwr.Message,
		) (response wwr.Payload, err error) {
			return nil, nil
		}
	}

	// Use default session manager if no specific one is defined
	if opts.SessionManager == nil {
		opts.SessionManager = newInMemSessManager()
	}

	// Use default address
	opts.Address = "127.0.0.1:0"

	// Use default heartbeat configuration if not set
	if opts.Heartbeat == wwr.OptionUnset {
		opts.Heartbeat = wwr.Disabled
	}

	server, err := wwr.NewServer(
		impl,
		opts,
	)
	if err != nil {
		t.Fatalf("Failed setting up server instance: %s", err)
	}

	// Run server in a separate goroutine
	go func() {
		if err := server.Run(); err != nil {
			panic(fmt.Errorf("Server failed: %s", err))
		}
	}()

	// Return reference to the server and the address its bound to
	return server
}

func comparePayload(t *testing.T, expected, actual wwr.Payload) {
	assert.Equal(t, expected.Encoding(), actual.Encoding())
	assert.Equal(t, expected.Data(), actual.Data())
}

func compareSessions(t *testing.T, expected, actual *wwr.Session) {
	if actual == nil && expected == nil {
		return
	}

	assert.NotNil(t, expected)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Key, actual.Key)
	assert.Equal(t, expected.Creation.Unix(), actual.Creation.Unix())
}
