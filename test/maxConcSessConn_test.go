package test

import (
	"context"
	"sync"
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
	sessionStorage := make(map[string]*wwr.Session)

	var sessionKey string
	sessionKeyLock := sync.RWMutex{}
	concurrentConns := uint(4)

	// Initialize server
	setup := setupTestServer(
		t,
		&serverImpl{
			onClientConnected: func(
				_ wwr.ConnectionOptions,
				conn wwr.Connection,
			) {
				// Created the session for the first connecting client only
				sessionKeyLock.Lock()
				defer sessionKeyLock.Unlock()
				if len(sessionKey) < 1 {
					assert.NoError(t, conn.CreateSession(nil))
					sessionKey = conn.SessionKey()
				}
			},
		},
		wwr.ServerOptions{
			MaxSessionConnections: concurrentConns,
			SessionManager: &callbackPoweredSessionManager{
				// Saves the session
				SessionCreated: func(conn wwr.Connection) error {
					sess := conn.Session()
					sessionStorage[sess.Key] = sess
					return nil
				},
				// Finds session by key
				SessionLookup: func(key string) (
					webwire.SessionLookupResult,
					error,
				) {
					if session, exists := sessionStorage[key]; exists {
						return webwire.NewSessionLookupResult(
							session.Creation,   // Creation
							session.LastLookup, // LastLookup
							webwire.SessionInfoToVarMap(
								session.Info,
							), // Info
						), nil
					}
					// Session not found
					return nil, nil
				},
			},
		},
		nil, // Use the default transport implementation
	)

	// Initialize client
	clients := make([]*testClient, concurrentConns)
	for i := uint(0); i < concurrentConns; i++ {
		client := setup.newClient(
			wwrclt.Options{
				DefaultRequestTimeout: 2 * time.Second,
			},
			nil, // Use the default transport implementation
			testClientHooks{},
		)
		clients[i] = client

		assert.NoError(t, client.connection.Connect())

		// Restore the session for all clients except the first one
		if i > 0 {
			sessionKeyLock.RLock()
			assert.NoError(t, client.connection.RestoreSession(
				context.Background(),
				[]byte(sessionKey),
			))
			sessionKeyLock.RUnlock()
		}
	}

	// Ensure that the last superfluous client is rejected
	superfluousClient := setup.newClient(
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		nil, // Use the default transport implementation
		testClientHooks{},
	)

	require.NoError(t, superfluousClient.connection.Connect())

	// Try to restore the session and expect this operation to fail
	// due to reached limit
	sessionKeyLock.RLock()
	sessionRestorationError := superfluousClient.connection.RestoreSession(
		context.Background(),
		[]byte(sessionKey),
	)
	require.Error(t, sessionRestorationError)
	require.IsType(t, wwr.MaxSessConnsReachedErr{}, sessionRestorationError)
	sessionKeyLock.RUnlock()
}
