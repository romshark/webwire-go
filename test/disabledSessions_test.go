package test

import (
	"context"
	"reflect"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestDisabledSessions tests errors returned by CreateSession, CloseSession
// and client.RestoreSession when sessions are disabled
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
		&serverImpl{
			onRequest: func(ctx context.Context) (wwr.Payload, error) {
				// Extract request message
				// and requesting client from the context
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
		wwr.ServerOptions{
			SessionsEnabled: false,
		},
	)

	// Initialize client
	client := wwrclt.NewClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
			Hooks: wwrclt.Hooks{
				OnSessionCreated: func(*wwr.Session) {
					t.Errorf("OnSessionCreated was not expected to be called")
				},
			},
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
