package test

import (
	"sync"
	"time"

	"github.com/qbeon/webwire-go"

	wwr "github.com/qbeon/webwire-go"
)

// inMemSessManager is a default in-memory session manager for testing purposes
type inMemSessManager struct {
	sessions map[string]*wwr.Client
	lock     sync.RWMutex
}

// newInMemSessManager constructs a new default session manager instance
// for testing purposes.
func newInMemSessManager() *inMemSessManager {
	return &inMemSessManager{
		sessions: make(map[string]*wwr.Client),
		lock:     sync.RWMutex{},
	}
}

// OnSessionCreated implements the session manager interface.
// It writes the created session into a file using the session key as file name
func (mng *inMemSessManager) OnSessionCreated(client *wwr.Client) error {
	mng.lock.Lock()
	mng.sessions[client.SessionKey()] = client
	mng.lock.Unlock()
	return nil
}

// OnSessionLookup implements the session manager interface.
// It searches the session file directory for the session file and loads it
func (mng *inMemSessManager) OnSessionLookup(key string) (
	bool,
	time.Time,
	map[string]interface{},
	error,
) {
	mng.lock.RLock()
	defer mng.lock.RUnlock()
	if clt, exists := mng.sessions[key]; exists {
		session := clt.Session()

		// Session found
		return true,
			session.Creation,
			webwire.SessionInfoToVarMap(session.Info),
			nil
	}

	// Session not found
	return false, time.Time{}, nil, nil
}

// OnSessionClosed implements the session manager interface.
// It closes the session by deleting the according session file
func (mng *inMemSessManager) OnSessionClosed(sessionKey string) error {
	mng.lock.Lock()
	delete(mng.sessions, sessionKey)
	mng.lock.Unlock()
	return nil
}

// callbackPoweredSessionManager represents a callback-powered session manager
// for testing purposes
type callbackPoweredSessionManager struct {
	SessionCreated func(client *wwr.Client) error
	SessionLookup  func(key string) (
		bool,
		time.Time,
		map[string]interface{},
		error,
	)
	SessionClosed func(sessionKey string) error
}

// OnSessionCreated implements the session manager interface
// calling the configured callback
func (mng *callbackPoweredSessionManager) OnSessionCreated(
	client *wwr.Client,
) error {
	if mng.SessionCreated == nil {
		return nil
	}
	return mng.SessionCreated(client)
}

// OnSessionLookup implements the session manager interface
// calling the configured callback
func (mng *callbackPoweredSessionManager) OnSessionLookup(
	key string,
) (bool, time.Time, map[string]interface{}, error) {
	if mng.SessionLookup == nil {
		return false, time.Time{}, nil, nil
	}
	return mng.SessionLookup(key)
}

// OnSessionClosed implements the session manager interface
// calling the configured callback
func (mng *callbackPoweredSessionManager) OnSessionClosed(
	sessionKey string,
) error {
	if mng.SessionClosed == nil {
		return nil
	}
	return mng.SessionClosed(sessionKey)
}
