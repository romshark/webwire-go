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
	_, addr := setupServer(
		t,
		wwr.ServerOptions{
			Hooks: wwr.Hooks{
				// Permanently store the session
				OnSessionCreated: func(_ *wwr.Client) error {
					return nil
				},
				// Find session by key
				OnSessionLookup: func(_ string) (*wwr.Session, error) {
					return nil, nil
				},
				// Define dummy hook to enable sessions on this server
				OnSessionClosed: func(_ *wwr.Client) error {
					return nil
				},
			},
		},
	)

	// Initialize client

	// Ensure that the last superfluous client is rejected
	client := wwrClient.NewClient(
		addr,
		wwrClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)

	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect client: %s", err)
	}

	// Try to restore the session and expect it to fail due to the session being inexistent
	sessRestErr := client.RestoreSession([]byte("lalala"))
	if _, isSessNotFoundErr := sessRestErr.(wwr.SessNotFound); !isSessNotFoundErr {
		t.Fatalf(
			"Expected a SessNotFound error during manual session restoration, got: %s | %s",
			reflect.TypeOf(sessRestErr),
			sessRestErr,
		)
	}
}
