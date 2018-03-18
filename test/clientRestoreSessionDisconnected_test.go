package test

import (
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientRestoreSessionDisconnected tests manual session restoration on disconnected client
func TestClientRestoreSessionDisconnected(t *testing.T) {
	// Initialize webwire server
	_, addr := setupServer(
		t,
		wwr.ServerOptions{},
	)

	// Initialize client
	client := wwrclt.NewClient(
		addr,
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)

	err := client.RestoreSession([]byte("somekey"))
	if _, isDisconnErr := err.(wwrclt.DisconnectedErr); !isDisconnErr {
		t.Fatalf(
			"Expected disconnected error, got: %s | %s",
			reflect.TypeOf(err),
			err,
		)
	}
}
