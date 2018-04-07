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
			onRequest: func(
				_ context.Context,
				clt *wwr.Client,
				_ *wwr.Message,
			) (wwr.Payload, error) {
				// Try to create a new session and expect an error
				createErr := clt.CreateSession(nil)
				verifyError(createErr)

				// Try to create a new session and expect an error
				closeErr := clt.CloseSession()
				verifyError(closeErr)

				return wwr.Payload{}, nil
			},
		},
		wwr.ServerOptions{
			SessionsEnabled: false,
		},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		func(*wwr.Session) {
			t.Errorf("OnSessionCreated was not expected to be called")
		},
		nil, nil, nil,
	)
	defer client.connection.Close()

	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send authentication request and await reply
	_, err := client.connection.Request(
		"login",
		wwr.Payload{Data: []byte("testdata")},
	)
	if err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	sessRestErr := client.connection.RestoreSession([]byte("testkey"))
	verifyError(sessRestErr)
}
