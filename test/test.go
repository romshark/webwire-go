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
	onClientConnected webwire.OnClientConnected,
	onClientDisconnected webwire.OnClientDisconnected,
	onSignal webwire.OnSignal,
	onRequest webwire.OnRequest,
	onSessionCreated webwire.OnSessionCreated,
	onSessionLookup webwire.OnSessionLookup,
	onSessionClosed webwire.OnSessionClosed,
	onOptions webwire.OnOptions,
) (srv *webwire.Server) {
	// Initialize webwire server
	webwireServer, err := webwire.NewServer(
		"127.0.0.1:0",
		onClientConnected,
		onClientDisconnected,
		onSignal, onRequest,
		onSessionCreated, onSessionLookup, onSessionClosed,
		onOptions,
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
	}
	// TODO: get session creation time from actual server time
	/*
		else if actual.CreationDate != expected.CreationDate {
			t.Errorf("Session creation dates differ:\n expected: '%s'\n actual:   '%s'",
				expected.CreationDate,
				actual.CreationDate,
			)
		}
	*/
	// TODO: implement user agent, OS and info comparison
	/*
		else if actual.UserAgent != expected.UserAgent {
			t.Errorf("Session user agents differ:\n expected: '%s'\n actual:   '%s'",
				expected.UserAgent,
				actual.UserAgent,
			)
		} else if actual.OperatingSystem != expected.OperatingSystem {
			t.Errorf("Session operating systems differ:\n expected: '%v'\n actual:   '%v'",
				expected.OperatingSystem,
				actual.OperatingSystem,
			)
		} else if !reflect.DeepEqual(expected.Info, actual.Info) {
			t.Errorf("Session info differs:\n expected: '%v'\n actual:   '%v'",
				expected.Info,
				actual.Info,
			)
		}
	*/
}
