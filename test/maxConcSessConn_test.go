package test

import (
	"reflect"
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
	server := setupServer(
		t,
		wwr.ServerOptions{
			SessionsEnabled:       true,
			MaxSessionConnections: concurrentConns,
			SessionManager: &CallbackPoweredSessionManager{
				// Saves the session
				SessionCreated: func(client *wwr.Client) error {
					sess := client.Session()
					sessionStorage[sess.Key] = sess
					return nil
				},
				// Finds session by key
				SessionLookup: func(key string) (*wwr.Session, error) {
					if session, exists := sessionStorage[key]; exists {
						return session, nil
					}
					return nil, nil
				},
			},
			Hooks: wwr.Hooks{
				OnClientConnected: func(client *wwr.Client) {
					// Created the session for the first connecting client only
					sessionKeyLock.Lock()
					defer sessionKeyLock.Unlock()
					if len(sessionKey) < 1 {
						if err := client.CreateSession(nil); err != nil {
							t.Errorf(
								"Unexpected error during session creation: %s",
								err,
							)
						}
						sessionKey = client.SessionKey()
					}
				},
			},
		},
	)

	// Initialize client
	clients := make([]*wwrClient.Client, concurrentConns)
	for i := uint(0); i < concurrentConns; i++ {
		client := wwrClient.NewClient(
			server.Addr().String(),
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
				t.Fatalf(
					"Unexpected error during manual session restoration: %s",
					err,
				)
			}
			sessionKeyLock.RUnlock()
		}
	}

	// Ensure that the last superfluous client is rejected
	superfluousClient := wwrClient.NewClient(
		server.Addr().String(),
		wwrClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)

	if err := superfluousClient.Connect(); err != nil {
		t.Fatalf("Couldn't connect superfluous client: %s", err)
	}

	// Try to restore the session and expect this operation to fail due to reached limit
	sessionKeyLock.RLock()
	sessRestErr := superfluousClient.RestoreSession([]byte(sessionKey))
	_, isMaxReachedErr := sessRestErr.(wwr.MaxSessConnsReachedErr)
	if !isMaxReachedErr {
		t.Fatalf(
			"Expected a MaxSessConnsReached error during "+
				"manual session restoration, got: %s | %s",
			reflect.TypeOf(sessRestErr),
			sessRestErr,
		)
	}
	sessionKeyLock.RUnlock()
}
