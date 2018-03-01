package test

import (
	"os"
	"reflect"
	"testing"

	webwire "github.com/qbeon/webwire-go"
)

// setupServer helps setting up
// and launching the server together with the hosting http server
func setupServer(
	t *testing.T,
	hooks webwire.Hooks,
) (srv *webwire.Server) {
	// Initialize webwire server
	webwireServer, err := webwire.NewServer(
		"127.0.0.1:0",
		hooks,
		os.Stdout, os.Stderr,
	)
	if err != nil {
		t.Fatalf("Failed creating a new WebWire server instance: %s", err)
	}

	return webwireServer
}

func comparePayload(t *testing.T, name string, expected, actual webwire.Payload) {
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

func compareSessions(t *testing.T, expected, actual *webwire.Session) {
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
