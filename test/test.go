package test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	wwr "github.com/qbeon/webwire-go"
)

// setupServer helps setting up
// and launching the server together with the hosting http server
func setupServer(t *testing.T, opts wwr.ServerOptions) (*wwr.Server, string) {
	// Setup headed server on arbitrary port
	opts.WarnLog = os.Stdout
	opts.ErrorLog = os.Stderr
	srv, _, addr, run, _, err := wwr.SetupServer(wwr.SetupOptions{
		ServerAddress: "127.0.0.1:0",
		ServerOptions: opts,
	})
	if err != nil {
		t.Fatalf("Failed setting up server instance: %s", err)
	}

	// Run server in a separate goroutine
	go func() {
		if err := run(); err != nil {
			panic(fmt.Errorf("Server failed: %s", err))
		}
	}()

	// Return reference to the server and the address its bound to
	return srv, addr
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
	} else if actual.CreationDate.Unix() != expected.CreationDate.Unix() {
		t.Errorf("Session creation dates differ:\n expected: '%s'\n actual:   '%s'",
			expected.CreationDate,
			actual.CreationDate,
		)
	}
}
