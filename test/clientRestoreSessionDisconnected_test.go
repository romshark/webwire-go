package test

import (
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestClientRestoreSessionDisconnected tests manual session restoration on disconnected client
// and expects client.RestoreSession to automatically establish a connection
func TestClientRestoreSessionDisconnected(t *testing.T) {
	// Initialize webwire server
	_, addr := setupServer(
		t,
		wwr.ServerOptions{
			SessionsEnabled: true,
			Hooks: wwr.Hooks{
				OnSessionCreated: func(_ *wwr.Client) error { return nil },
				OnSessionLookup:  func(_ string) (*wwr.Session, error) { return nil, nil },
				OnSessionClosed:  func(_ *wwr.Client) error { return nil },
			},
		},
	)

	// Initialize client and skip manual connection establishment
	client := wwrclt.NewClient(
		addr,
		wwrclt.Options{
			DefaultRequestTimeout: 100 * time.Millisecond,
		},
	)

	err := client.RestoreSession([]byte("inexistentkey"))
	if _, isSessNotFoundErr := err.(wwr.SessNotFoundErr); !isSessNotFoundErr {
		t.Fatalf(
			"Expected disconnected error, got: %s | %s",
			reflect.TypeOf(err),
			err,
		)
	}
}
