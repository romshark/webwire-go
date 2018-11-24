package webwire

import (
	"fmt"
	"sync"
)

// sessionRegistry represents a thread safe registry
// of all currently active sessions
type sessionRegistry struct {
	server           *server
	lock             *sync.RWMutex
	maxConns         uint
	registry         map[string]map[*connection]struct{}
	onSessionDestroy func(sessionKey string)
}

// newSessionRegistry returns a new instance of a session registry.
// maxConns defines the maximum number of concurrent connections
// for a single session while zero stands for unlimited
func newSessionRegistry(
	maxConns uint,
	onSessionDestroy func(sessionKey string),
) *sessionRegistry {
	return &sessionRegistry{
		lock:             &sync.RWMutex{},
		maxConns:         maxConns,
		registry:         make(map[string]map[*connection]struct{}),
		onSessionDestroy: onSessionDestroy,
	}
}

// register registers a new connection for the given clients session.
// Returns an error if the given clients session already reached
// the maximum number of concurrent connections
func (asr *sessionRegistry) register(con *connection) error {
	asr.lock.Lock()
	if connSet, exists := asr.registry[con.session.Key]; exists {
		// Ensure max connections isn't exceeded
		if asr.maxConns > 0 && uint(len(connSet)+1) > asr.maxConns {
			asr.lock.Unlock()
			return fmt.Errorf(
				"Max conns (%d) reached for session %s",
				asr.maxConns,
				con.session.Key,
			)
		}
		// Overwrite the current entry incrementing the number of connections
		connSet[con] = struct{}{}
		asr.registry[con.session.Key] = connSet
		asr.lock.Unlock()
		return nil
	}
	newList := map[*connection]struct{}{
		con: {},
	}
	asr.registry[con.session.Key] = newList
	asr.lock.Unlock()
	return nil
}

// deregister removes a connection from the list of connections of a session
// and returns the number of connections left.
// If there's only one connection left then the entire session will be removed
// from the register and 0 will be returned. If destroy is set to true then the
// session manager hook OnSessionClosed is executed destroying the session,
// otherwise the session will be deregistered but left untouched in the storage.
// If the given connection is not in the register -1 is returned
func (asr *sessionRegistry) deregister(
	conn *connection,
	destroy bool,
) int {
	if conn.session == nil {
		return -1
	}

	asr.lock.Lock()
	if connSet, exists := asr.registry[conn.session.Key]; exists {
		// If a single connection is left then remove or destroy the session
		if len(connSet) < 2 {
			delete(asr.registry, conn.session.Key)
			asr.lock.Unlock()

			// Destroy the session
			if destroy {
				// Recover potential user-space hook panics
				// to avoid panicking the server
				defer func() {
					if recvErr := recover(); recvErr != nil {
						asr.server.errorLog.Printf(
							"session closure hook panic: %v",
							recvErr,
						)
					}
				}()

				asr.onSessionDestroy(conn.session.Key)
			}

			return 0
		}

		// Find and remove the client from the connections list
		delete(connSet, conn)
		connSetLen := len(connSet)
		asr.lock.Unlock()
		return connSetLen
	}

	asr.lock.Unlock()
	return -1
}

// activeSessionsNum returns the number of currently active sessions
func (asr *sessionRegistry) activeSessionsNum() int {
	asr.lock.RLock()
	registryLen := len(asr.registry)
	asr.lock.RUnlock()
	return registryLen
}

// sessionConnectionsNum implements the sessionRegistry interface
func (asr *sessionRegistry) sessionConnectionsNum(sessionKey string) int {
	asr.lock.RLock()
	if connSet, exists := asr.registry[sessionKey]; exists {
		connSetLen := len(connSet)
		asr.lock.RUnlock()
		return connSetLen
	}
	asr.lock.RUnlock()
	return -1
}

// sessionConnections implements the sessionRegistry interface
func (asr *sessionRegistry) sessionConnections(
	sessionKey string,
) map[*connection]struct{} {
	asr.lock.RLock()
	if connSet, exists := asr.registry[sessionKey]; exists {
		asr.lock.RUnlock()
		return connSet
	}
	asr.lock.RUnlock()
	return nil
}
