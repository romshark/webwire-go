package test

import (
	"sync"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	wwrClient "github.com/qbeon/webwire-go/client"
)

// TestMaxConcSessConn tests 4 maximum concurrent connections of a session
func TestMaxConcSessConn(t *testing.T) {
	sessionStorage := make(map[string]*wwr.Session)

	var sessionKey string
	sessionKeyLock := sync.RWMutex{}
	concurrentConns := uint(4)

	// Initialize server
	_, addr := setupServer(
		t,
		wwr.ServerOptions{
			MaxSessionConnections: concurrentConns,
			Hooks: wwr.Hooks{
				OnClientConnected: func(client *wwr.Client) {
					// Created the session for the first connecting client only
					sessionKeyLock.Lock()
					defer sessionKeyLock.Unlock()
					if len(sessionKey) < 1 {
						if err := client.CreateSession(nil); err != nil {
							t.Errorf("Unexpected error during session creation: %s", err)
						}
						sessionKey = client.Session.Key
					}
				},
				// Permanently store the session
				OnSessionCreated: func(client *wwr.Client) error {
					sessionStorage[client.Session.Key] = client.Session
					return nil
				},
				// Find session by key
				OnSessionLookup: func(key string) (*wwr.Session, error) {
					if session, exists := sessionStorage[key]; exists {
						return session, nil
					}
					return nil, nil
				},
				// Define dummy hook to enable sessions on this server
				OnSessionClosed: func(_ *wwr.Client) error { return nil },
			},
		},
	)

	// Initialize client
	clients := make([]*wwrClient.Client, concurrentConns)
	for i := uint(0); i < concurrentConns; i++ {
		client := wwrClient.NewClient(
			addr,
			wwrClient.Options{
				DefaultRequestTimeout: 2 * time.Second,
			},
		)
		clients[i] = client

		if err := client.Connect(); err != nil {
			t.Fatalf("Couldn't connect client: %s", err)
		}

		// Restore the session for all clients except the first one
		if i > 0 {
			sessionKeyLock.RLock()
			if err := client.RestoreSession([]byte(sessionKey)); err != nil {
				t.Fatalf("Unexpected error during manual session restoration: %s", err)
			}
			sessionKeyLock.RUnlock()
		}
	}

	// Ensure that the last superfluous client is rejected
	superflousClient := wwrClient.NewClient(
		addr,
		wwrClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)

	if err := superflousClient.Connect(); err != nil {
		t.Fatalf("Couldn't connect superfluous client: %s", err)
	}

	// Try to restore the session and expect this operation to fail due to reached limit
	sessionKeyLock.RLock()
	if err := superflousClient.RestoreSession([]byte(sessionKey)); err == nil {
		t.Fatalf("Expected an error during superfluous client manual session restoration")
	}
	sessionKeyLock.RUnlock()
}
