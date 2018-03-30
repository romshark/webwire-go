package test

import (
	"context"
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestDisabledSessions verifies the server is connectable,
// and is able to receives requests and signals, create sessions
// and identify clients during request- and signal handling
func TestDisabledSessions(t *testing.T) {
	verifyError := func(err error) {
		if _, isDisabledErr := err.(wwr.SessionsDisabledErr); !isDisabledErr {
			t.Errorf(
				"Expected SessionsDisabled error, got: %s | %s",
				reflect.TypeOf(err),
				err,
			)
		}
	}

	// Initialize webwire server
	server := setupServer(
		t,
		wwr.ServerOptions{
			SessionsEnabled: false,
			Hooks: wwr.Hooks{
				OnRequest: func(ctx context.Context) (wwr.Payload, error) {
					// Extract request message and requesting client from the context
					msg := ctx.Value(wwr.Msg).(wwr.Message)

					// Try to create a new session and expect an error
					createErr := msg.Client.CreateSession(nil)
					verifyError(createErr)

					// Try to create a new session and expect an error
					closeErr := msg.Client.CloseSession()
					verifyError(closeErr)

					return wwr.Payload{}, nil
				},
			},
		},
	)

	// Initialize client
	client := wwrclt.NewClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)
	defer client.Close()

	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send authentication request and await reply
	_, err := client.Request("login", wwr.Payload{Data: []byte("testdata")})
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	sessRestErr := client.RestoreSession([]byte("testkey"))
	verifyError(sessRestErr)
}
