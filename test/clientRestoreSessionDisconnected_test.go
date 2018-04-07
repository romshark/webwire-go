package test

import (
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientRestoreSessionDisconnected tests manual session restoration
// on disconnected client expecting client.RestoreSession
// to automatically establish a connection
func TestClientRestoreSessionDisconnected(t *testing.T) {
	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{
			SessionsEnabled: true,
		},
	)

	// Initialize client and skip manual connection establishment
	client := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 100 * time.Millisecond,
		},
		nil, nil, nil, nil,
	)

	err := client.connection.RestoreSession([]byte("inexistentkey"))
	if _, isSessNotFoundErr := err.(wwr.SessNotFoundErr); !isSessNotFoundErr {
		t.Fatalf(
			"Expected disconnected error, got: %s | %s",
			reflect.TypeOf(err),
			err,
		)
	}
}
