package test

import (
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientRestSessAutoconn tests manual session restoration
// on disconnected client expecting client.RestoreSession
// to automatically establish a connection
func TestClientRestSessAutoconn(t *testing.T) {
	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 100 * time.Millisecond,
		},
		callbackPoweredClientHooks{},
	)

	// Skip manual connection establishment and rely on autoconnect instead
	err := client.connection.RestoreSession([]byte("inexistentkey"))
	if _, isSessNotFoundErr := err.(wwr.SessNotFoundErr); !isSessNotFoundErr {
		t.Fatalf(
			"Expected disconnected error, got: %s | %s",
			reflect.TypeOf(err),
			err,
		)
	}
}
