package test

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	wwr "github.com/qbeon/webwire-go"
)

// setupServer helps setting up and launching the server
// together with the hosting http server
// setting up a headed server on a randomly assigned port
func setupServer(
	t *testing.T,
	impl *serverImpl,
	opts wwr.ServerOptions,
) wwr.WebwireServer {
	// Setup headed server on arbitrary port

	if impl.beforeUpgrade == nil {
		impl.beforeUpgrade = func(_ http.ResponseWriter, _ *http.Request) bool {
			return true
		}
	}
	if impl.onClientConnected == nil {
		impl.onClientConnected = func(_ *wwr.Client) {}
	}
	if impl.onClientDisconnected == nil {
		impl.onClientDisconnected = func(_ *wwr.Client) {}
	}
	if impl.onSignal == nil {
		impl.onSignal = func(
			_ context.Context,
			_ *wwr.Client,
			_ *wwr.Message,
		) {
		}
	}
	if impl.onRequest == nil {
		impl.onRequest = func(
			_ context.Context,
			_ *wwr.Client,
			_ *wwr.Message,
		) (response wwr.Payload, err error) {
			return wwr.Payload{}, nil
		}
	}

	// Use default session manager if no specific one is defined
	if opts.SessionManager == nil {
		opts.SessionManager = newInMemSessManager()
	}

	opts.Address = "127.0.0.1:0"

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

func comparePayload(t *testing.T, name string, expected, actual wwr.Payload) {
	if actual.Encoding != expected.Encoding {
		t.Errorf(
			"Invalid %s: payload encoding differs:"+
				"\n expected: '%v'\n actual:   '%v'",
			name,
			expected.Encoding,
			actual.Encoding,
		)
		return
	}
	if !reflect.DeepEqual(actual.Data, expected.Data) {
		t.Errorf(
			"Invalid %s: payload data differs:"+
				"\n expected: '%s'\n actual:   '%s'",
			name,
			string(expected.Data),
			string(actual.Data),
		)
	}
}

func compareSessions(t *testing.T, expected, actual *wwr.Session) {
	if actual == nil && expected == nil {
		return
	} else if (actual == nil && expected != nil) ||
		(expected == nil && actual != nil) {
		t.Errorf("Sessions differ:\n expected: '%v'\n actual:   '%v'",
			expected,
			actual,
		)
	} else if actual.Key != expected.Key {
		t.Errorf("Session keys differ:\n expected: '%s'\n actual:   '%s'",
			expected.Key,
			actual.Key,
		)
	} else if actual.Creation.Unix() != expected.Creation.Unix() {
		t.Errorf("Session creation dates differ:\n"+
			" expected: '%s'\n actual:   '%s'",
			expected.Creation,
			actual.Creation,
		)
	}
}
