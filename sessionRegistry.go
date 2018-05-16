package webwire

import (
	"fmt"
	"sync"
)

// sessionRegistry represents a thread safe registry of all currently active sessions
type sessionRegistry struct {
	lock     sync.RWMutex
	maxConns uint
	registry map[string][]*Client
}

// newSessionRegistry returns a new instance of a session registry.
// maxConns defines the maximum number of concurrent connections for a single session
// while zero stands for unlimited
func newSessionRegistry(maxConns uint) *sessionRegistry {
	return &sessionRegistry{
		lock:     sync.RWMutex{},
		maxConns: maxConns,
		registry: make(map[string][]*Client),
	}
}

// register registers a new connection for the given clients session.
// Returns an error if the given clients session already reached
// the maximum number of concurrent connections
func (asr *sessionRegistry) register(clt *Client) error {
	asr.lock.Lock()
	defer asr.lock.Unlock()
	if agentList, exists := asr.registry[clt.session.Key]; exists {
		// Ensure max connections isn't exceeded
		if asr.maxConns > 0 && uint(len(agentList)+1) > asr.maxConns {
			return fmt.Errorf(
				"Max conns (%d) reached for session %s",
				asr.maxConns,
				clt.session.Key,
			)
		}
		// Overwrite the current entry incrementing the number of connections
		asr.registry[clt.session.Key] = append(agentList, clt)
		return nil
	}
	newList := []*Client{clt}
	asr.registry[clt.session.Key] = newList
	return nil
}

// deregister removes a client agent from the list of connections of a session
// returns the number of connections left.
// If there's only one connection left then the entire session will be removed
// from the register and 0 will be returned.
// If the given client agent is not in the register -1 is returned
func (asr *sessionRegistry) deregister(clt *Client) int {
	if clt.session == nil {
		return -1
	}

	asr.lock.Lock()
	defer asr.lock.Unlock()
	if agentList, exists := asr.registry[clt.session.Key]; exists {
		// If a single connection is left then remove the session
		if len(agentList) < 2 {
			delete(asr.registry, clt.session.Key)
			return 0
		}
		// Find and remove the client from the connections list
		for index, agent := range agentList {
			if agent == clt {
				asr.registry[clt.session.Key] = append(
					agentList[:index],
					agentList[index+1:]...,
				)
			}
		}
		return len(agentList) - 1
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
	defer asr.lock.RUnlock()
	if agentList, exists := asr.registry[sessionKey]; exists {
		return len(agentList)
	}
	return -1
}

// sessionConnections implements the sessionRegistry interface
func (asr *sessionRegistry) sessionConnections(sessionKey string) []*Client {
	asr.lock.RLock()
	defer asr.lock.RUnlock()
	if agentList, exists := asr.registry[sessionKey]; exists {
		return agentList
	}
	return nil
}
