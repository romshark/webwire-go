package webwire

import (
	"sync"
)

// sessionRegistryEntry represents a session registry entry
type sessionRegistryEntry struct {
	connections uint
	client      *Client
}

// sessionRegistry represents a thread safe registry of all currently active sessions
type sessionRegistry struct {
	lock     sync.RWMutex
	maxConns uint
	registry map[string]sessionRegistryEntry
}

// newSessionRegistry returns a new instance of a session registry.
// maxConns defines the maximum number of concurrent connections for a single session
// while zero stands for unlimited
func newSessionRegistry(maxConns uint) sessionRegistry {
	return sessionRegistry{
		lock:     sync.RWMutex{},
		maxConns: maxConns,
		registry: make(map[string]sessionRegistryEntry),
	}
}

// register registers a new connection for the given clients session and returns true.
// Returns false if the given clients session already has the max number of connections assigned.
func (asr *sessionRegistry) register(clt *Client) bool {
	asr.lock.Lock()
	defer asr.lock.Unlock()
	if entry, exists := asr.registry[clt.Session.Key]; exists {
		// Ensure max connections isn't exceeded
		if asr.maxConns > 0 && entry.connections+1 > asr.maxConns {
			return false
		}
		// Overwrite the current entry incrementing the number of connections
		asr.registry[clt.Session.Key] = sessionRegistryEntry{
			connections: entry.connections + 1,
			client:      entry.client,
		}
		return true
	}
	asr.registry[clt.Session.Key] = sessionRegistryEntry{
		connections: 1,
		client:      clt,
	}
	return true
}

// deregister decrements the number of connections assigned to the given clients session
// and returns true. If there's only one connection left then the session will be removed
// from the register and false will be returned
func (asr *sessionRegistry) deregister(clt *Client) bool {
	asr.lock.Lock()
	defer asr.lock.Unlock()
	if entry, exists := asr.registry[clt.Session.Key]; exists {
		// If a single connection is left then remove the session
		if entry.connections < 2 {
			delete(asr.registry, clt.Session.Key)
			return false
		}
		// Overwrite the current entry decrementing the number of connections
		asr.registry[clt.Session.Key] = sessionRegistryEntry{
			connections: entry.connections - 1,
			client:      entry.client,
		}
	}
	return false
}

// ActiveSessions returns the number of currently active sessions
func (asr *sessionRegistry) ActiveSessions() int {
	asr.lock.RLock()
	len := len(asr.registry)
	asr.lock.RUnlock()
	return len
}

// SessionConnections returns the number of concurrent connections
// associated with the session associated with the given key.
// Returns zero if the session associated with the given key doesn't exist.
func (asr *sessionRegistry) SessionConnections(sessionKey string) uint {
	asr.lock.RLock()
	defer asr.lock.RUnlock()
	if sess, exists := asr.registry[sessionKey]; exists {
		return sess.connections
	}
	return 0
}
