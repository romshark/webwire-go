package test

import (
	"sync"
	"time"

	wwr "github.com/qbeon/webwire-go"
)

type session struct {
	Key        string
	Creation   time.Time
	LastLookup time.Time
	Info       wwr.SessionInfo
}

// inMemSessManager is a default in-memory session manager for testing purposes
type inMemSessManager struct {
	sessions map[string]session
	lock     sync.Mutex
}

// newInMemSessManager constructs a new default session manager instance
// for testing purposes.
func newInMemSessManager() *inMemSessManager {
	return &inMemSessManager{
		sessions: make(map[string]session),
		lock:     sync.Mutex{},
	}
}

// OnSessionCreated implements the session manager interface.
// It writes the created session into a file using the session key as file name
func (mng *inMemSessManager) OnSessionCreated(conn wwr.Connection) error {
	mng.lock.Lock()
	sess := conn.Session()
	var sessInfo wwr.SessionInfo
	if sess.Info != nil {
		sessInfo = sess.Info.Copy()
	}
	mng.sessions[sess.Key] = session{
		Key:      sess.Key,
		Creation: sess.Creation,
		Info:     sessInfo,
	}
	mng.lock.Unlock()
	return nil
}

// OnSessionLookup implements the session manager interface.
// It searches the session file directory for the session file and loads it
func (mng *inMemSessManager) OnSessionLookup(key string) (
	wwr.SessionLookupResult,
	error,
) {
	mng.lock.Lock()
	defer mng.lock.Unlock()
	if session, exists := mng.sessions[key]; exists {
		// Update last lookup field
		session.LastLookup = time.Now().UTC()
		mng.sessions[key] = session

		// Session found
		return wwr.SessionLookupResult{
			Creation:   session.Creation,
			LastLookup: session.LastLookup,
			Info:       wwr.SessionInfoToVarMap(session.Info),
		}, nil
	}

	// Session not found
	return wwr.SessionLookupResult{}, wwr.SessNotFoundErr{}
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
	SessionCreated func(client wwr.Connection) error
	SessionLookup  func(key string) (
		wwr.SessionLookupResult,
		error,
	)
	SessionClosed func(sessionKey string) error
}

// OnSessionCreated implements the session manager interface
// calling the configured callback
func (mng *callbackPoweredSessionManager) OnSessionCreated(
	client wwr.Connection,
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
) (wwr.SessionLookupResult, error) {
	if mng.SessionLookup == nil {
		return wwr.SessionLookupResult{}, wwr.SessNotFoundErr{}
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
