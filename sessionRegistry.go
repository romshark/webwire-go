package webwire

import (
	"sync"
)

// sessionRegistry represents a thread safe registry of all currently active sessions
type sessionRegistry struct {
	lock     sync.RWMutex
	registry map[string]*Client
}

// newSessionRegistry returns a new instance of a session registry
func newSessionRegistry() sessionRegistry {
	return sessionRegistry{
		lock:     sync.RWMutex{},
		registry: make(map[string]*Client),
	}
}

// register registers the given clients session as a currently active session
func (asr *sessionRegistry) register(clt *Client) {
	asr.lock.Lock()
	asr.registry[clt.Session.Key] = clt
	asr.lock.Unlock()
}

// deregister deregisters the given clients session from the list of currently active sessions
func (asr *sessionRegistry) deregister(clt *Client) {
	asr.lock.Lock()
	delete(asr.registry, clt.Session.Key)
	asr.lock.Unlock()
}

// Len returns the number of currently active sessions
func (asr *sessionRegistry) Len() int {
	asr.lock.RLock()
	len := len(asr.registry)
	asr.lock.RUnlock()
	return len
}

// Exists returns true if the session associated with the given key exists and is currently active
func (asr *sessionRegistry) Exists(sessionKey string) bool {
	asr.lock.RLock()
	_, exists := asr.registry[sessionKey]
	asr.lock.RUnlock()
	return exists
}
