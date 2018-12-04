package test

import (
	"context"
	"testing"
	"time"

	"github.com/qbeon/webwire-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMaxConcSessConn tests 4 maximum concurrent connections of a session
func TestMaxConcSessConn(t *testing.T) {
	concurrentConns := uint(4)

	var sessionKey = "testsessionkey"
	sessionCreation := time.Now()

	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{},
		wwr.ServerOptions{
			MaxSessionConnections: concurrentConns,
			SessionManager: &SessionManager{
				SessionLookup: func(key string) (
					webwire.SessionLookupResult,
					error,
				) {
					if key != sessionKey {
						// Session not found
						return nil, nil
					}
					return webwire.NewSessionLookupResult(
						sessionCreation, // Creation
						time.Now(),      // LastLookup
						nil,             // Info
					), nil
				},
			},
		},
		nil, // Use the default transport implementation
	)

	clientOptions := wwrclt.Options{
		DefaultRequestTimeout: 2 * time.Second,
	}

	// Initialize clients
	clients := make([]*TestClient, concurrentConns)
	for i := uint(0); i < concurrentConns; i++ {
		client := setup.NewClient(
			clientOptions,
			nil, // Use the default transport implementation
			TestClientHooks{},
		)
		clients[i] = client

		// Restore the session for all clients
		assert.NoError(t, client.Connection.RestoreSession(
			context.Background(),
			[]byte(sessionKey),
		))
	}

	// Ensure that the last superfluous client is rejected
	superfluousClient := setup.NewClient(
		clientOptions,
		nil, // Use the default transport implementation
		TestClientHooks{},
	)

	require.NoError(t, superfluousClient.Connection.Connect())

	// Try to restore the session and expect this operation to fail
	// due to reached limit
	sessionRestorationError := superfluousClient.Connection.RestoreSession(
		context.Background(),
		[]byte(sessionKey),
	)
	require.Error(t, sessionRestorationError)
	require.IsType(t, wwr.MaxSessConnsReachedErr{}, sessionRestorationError)
}
