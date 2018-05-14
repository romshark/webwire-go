package test

import (
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrClient "github.com/qbeon/webwire-go/client"
)

// TestRestoreInexistentSession tests the restoration of an inexistent session
func TestRestoreInexistentSession(t *testing.T) {
	// Initialize server
	server := setupServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{},
	)

	// Initialize client

	// Ensure that the last superfluous client is rejected
	client := newCallbackPoweredClient(
		server.Addr().String(),
		wwrClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect client: %s", err)
	}

	// Try to restore the session and expect it to fail
	// due to the session being inexistent
	sessRestErr := client.connection.RestoreSession([]byte("lalala"))
	_, isSessNotFoundErr := sessRestErr.(wwr.SessNotFoundErr)
	if !isSessNotFoundErr {
		t.Fatalf(
			"Expected a SessNotFound error during manual session restoration"+
				", got: %s | %s",
			reflect.TypeOf(sessRestErr),
			sessRestErr,
		)
	}
}
