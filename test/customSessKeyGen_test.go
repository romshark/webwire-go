package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

// TestCustomSessKeyGen tests custom session key generators
func TestCustomSessKeyGen(t *testing.T) {
	expectedSessionKey := "customkey123"

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Try to create a new session
				err := conn.CreateSession(nil)
				assert.NoError(t, err)
				if err != nil {
					return nil, err
				}

				key := conn.SessionKey()
				assert.Equal(t, expectedSessionKey, key)

				// Return the key of the newly created session
				// (use default binary encoding)
				return wwr.NewPayload(
					wwr.EncodingBinary,
					[]byte(key),
				), nil
			},
		},
		wwr.ServerOptions{
			SessionKeyGenerator: &sessionKeyGen{
				generate: func() string {
					return expectedSessionKey
				},
			},
		},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{},
	)
	defer client.connection.Close()

	// Send authentication request and await reply
	_, err := client.connection.Request(
		context.Background(),
		[]byte("login"),
		wwr.NewPayload(wwr.EncodingBinary, []byte("testdata")),
	)
	require.NoError(t, err)

	require.Equal(t, expectedSessionKey, client.connection.Session().Key)
}
