package test

import (
	"sync"

	wwr "github.com/qbeon/webwire-go"
)

// InMemSessManager is a default in-memory session manager for testing purposes
type InMemSessManager struct {
	sessions map[string]*wwr.Client
	lock     sync.RWMutex
}

// NewInMemSessManager constructs a new default session manager instance
// for testing purposes.
func NewInMemSessManager() *InMemSessManager {
	return &InMemSessManager{
		sessions: make(map[string]*wwr.Client),
		lock:     sync.RWMutex{},
	}
}

// OnSessionCreated implements the session manager interface.
// It writes the created session into a file using the session key as file name
func (mng *InMemSessManager) OnSessionCreated(client *wwr.Client) error {
	mng.lock.Lock()
	mng.sessions[client.SessionKey()] = client
	mng.lock.Unlock()
	return nil
}

// OnSessionLookup implements the session manager interface.
// It searches the session file directory for the session file and loads it
func (mng *InMemSessManager) OnSessionLookup(key string) (*wwr.Session, error) {
	mng.lock.RLock()
	defer mng.lock.RUnlock()
	if clt, exists := mng.sessions[key]; exists {
		return clt.Session(), nil
	}
	return nil, nil
}

// OnSessionClosed implements the session manager interface.
// It closes the session by deleting the according session file
func (mng *InMemSessManager) OnSessionClosed(client *wwr.Client) error {
	mng.lock.Lock()
	delete(mng.sessions, client.SessionKey())
	mng.lock.Unlock()
	return nil
}

// CallbackPoweredSessionManager represents a callback-powered session manager
// for testing purposes
type CallbackPoweredSessionManager struct {
	SessionCreated func(client *wwr.Client) error
	SessionLookup  func(key string) (*wwr.Session, error)
	SessionClosed  func(client *wwr.Client) error
}

// OnSessionCreated implements the session manager interface
// calling the configured callback
func (mng *CallbackPoweredSessionManager) OnSessionCreated(
	client *wwr.Client,
) error {
	if mng.SessionCreated == nil {
		return nil
	}
	return mng.SessionCreated(client)
}

// OnSessionLookup implements the session manager interface
// calling the configured callback
func (mng *CallbackPoweredSessionManager) OnSessionLookup(
	key string,
) (*wwr.Session, error) {
	if mng.SessionLookup == nil {
		return nil, nil
	}
	return mng.SessionLookup(key)
}

// OnSessionClosed implements the session manager interface
// calling the configured callback
func (mng *CallbackPoweredSessionManager) OnSessionClosed(
	client *wwr.Client,
) error {
	if mng.SessionClosed == nil {
		return nil
	}
	return mng.SessionClosed(client)
}
