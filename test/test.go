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

func comparePayload(t *testing.T, name string, expected, actual []byte) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(
			"Invalid %s: payload doesn't match:"+
				"\n expected: '%s'\n actual:   '%s'",
			name,
			string(expected),
			string(actual),
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
