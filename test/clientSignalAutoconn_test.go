package test

import (
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientSignalAutoconn tests sending signals on disconnected clients
// expecting it to automatically connect
func TestClientSignalAutoconn(t *testing.T) {
	// Initialize webwire server given only the request
	server := setupServer(
		t,
		&serverImpl{},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)

	// Skip manual connection establishment and rely on autoconnect instead
	if err := client.connection.Signal(
		"",
		wwr.NewPayload(wwr.EncodingBinary, []byte("testdata")),
	); err != nil {
		t.Fatalf(
			"Expected signal to automatically connect, got error: %s | %s",
			reflect.TypeOf(err),
			err,
		)
	}
}
