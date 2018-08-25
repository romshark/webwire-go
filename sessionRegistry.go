package webwire

import (
	"fmt"
	"sync"
)

// sessionRegistry represents a thread safe registry
// of all currently active sessions
type sessionRegistry struct {
	lock     sync.RWMutex
	maxConns uint
	registry map[string][]*connection
}

// newSessionRegistry returns a new instance of a session registry.
// maxConns defines the maximum number of concurrent connections
// for a single session while zero stands for unlimited
func newSessionRegistry(maxConns uint) *sessionRegistry {
	return &sessionRegistry{
		lock:     sync.RWMutex{},
		maxConns: maxConns,
		registry: make(map[string][]*connection),
	}
}

// register registers a new connection for the given clients session.
// Returns an error if the given clients session already reached
// the maximum number of concurrent connections
func (asr *sessionRegistry) register(con *connection) error {
	asr.lock.Lock()
	defer asr.lock.Unlock()
	if connList, exists := asr.registry[con.session.Key]; exists {
		// Ensure max connections isn't exceeded
		if asr.maxConns > 0 && uint(len(connList)+1) > asr.maxConns {
			return fmt.Errorf(
				"Max conns (%d) reached for session %s",
				asr.maxConns,
				con.session.Key,
			)
		}
		// Overwrite the current entry incrementing the number of connections
		asr.registry[con.session.Key] = append(connList, con)
		return nil
	}
	newList := []*connection{con}
	asr.registry[con.session.Key] = newList
	return nil
}

// deregister removes a connection from the list of connections of a session
// returns the number of connections left.
// If there's only one connection left then the entire session will be removed
// from the register and 0 will be returned.
// If the given connection is not in the register -1 is returned
func (asr *sessionRegistry) deregister(con *connection) int {
	if con.session == nil {
		return -1
	}

	asr.lock.Lock()
	defer asr.lock.Unlock()
	if connList, exists := asr.registry[con.session.Key]; exists {
		// If a single connection is left then remove the session
		if len(connList) < 2 {
			delete(asr.registry, con.session.Key)
			return 0
		}
		// Find and remove the client from the connections list
		for index, conn := range connList {
			if conn == con {
				asr.registry[con.session.Key] = append(
					connList[:index],
					connList[index+1:]...,
				)
			}
		}
		return len(connList) - 1
	}
	return -1
}

// activeSessionsNum returns the number of currently active sessions
func (asr *sessionRegistry) activeSessionsNum() int {
	asr.lock.RLock()
	len := len(asr.registry)
	asr.lock.RUnlock()
	return len
}

// sessionConnectionsNum implements the sessionRegistry interface
func (asr *sessionRegistry) sessionConnectionsNum(sessionKey string) int {
	asr.lock.RLock()
	if connList, exists := asr.registry[sessionKey]; exists {
		len := len(connList)
		asr.lock.RUnlock()
		return len
	}
	asr.lock.RUnlock()
	return -1
}

// sessionConnections implements the sessionRegistry interface
func (asr *sessionRegistry) sessionConnections(
	sessionKey string,
) []*connection {
	asr.lock.RLock()
	if connList, exists := asr.registry[sessionKey]; exists {
		asr.lock.RUnlock()
		return connList
	}
	asr.lock.RUnlock()
	return nil
}
