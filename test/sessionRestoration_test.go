package test

import (
	"sync"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
)

// TestSessionRestoration tests manual session restoration by key
func TestSessionRestoration(t *testing.T) {
	lookupTriggered := sync.WaitGroup{}
	lookupTriggered.Add(1)
	var sessionKey = "testsessionkey"
	sessionCreation := time.Now()

	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{},
		wwr.ServerOptions{
			SessionManager: &SessionManager{
				SessionLookup: func(key string) (
					wwr.SessionLookupResult,
					error,
				) {
					defer lookupTriggered.Done()
					if key != sessionKey {
						// Session not found
						return nil, nil
					}
					return wwr.NewSessionLookupResult(
						sessionCreation, // Creation
						time.Now(),      // LastLookup
						nil,             // Info
					), nil
				},
			},
		},
		nil, // Use the default transport implementation
	)

	// Initialize clients
	sock, _ := setup.NewClientSocket()

	requestRestoreSessionSuccess(t, sock, []byte(sessionKey))

	lookupTriggered.Wait()
}
