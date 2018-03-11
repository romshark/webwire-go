package state

import (
	"sync"

	wwr "github.com/qbeon/webwire-go"
)

type inMemoryStorage struct {
	lock      sync.RWMutex
	connected map[*wwr.Client]bool
	sessions  map[string]*wwr.Client
}

// AddConnected appends the given client to the list of connected clients
func (strg *inMemoryStorage) AddConnected(client *wwr.Client) {
	strg.lock.Lock()
	defer strg.lock.Unlock()
	strg.connected[client] = true
}

// RemoveConnected removes the given client from the list of connected clients
func (strg *inMemoryStorage) RemoveConnected(client *wwr.Client) {
	strg.lock.Lock()
	defer strg.lock.Unlock()
	delete(strg.connected, client)
}

// NumConnected returns the number of connected clients
func (strg *inMemoryStorage) NumConnected() int {
	strg.lock.RLock()
	defer strg.lock.RUnlock()
	return len(strg.connected)
}

func (strg *inMemoryStorage) ForEachConnected(lambda func(*wwr.Client)) {
	strg.lock.RLock()
	defer strg.lock.RUnlock()
	for clt := range strg.connected {
		lambda(clt)
	}
}

// SaveSession saves the given clients session in memory
func (strg *inMemoryStorage) SaveSession(client *wwr.Client) {
	strg.lock.Lock()
	defer strg.lock.Unlock()
	strg.sessions[client.Session.Key] = client
}

// FindSession searches for the session by key,
// returns nil if there's none associated the given key
func (strg *inMemoryStorage) FindSession(key string) (*wwr.Session, error) {
	strg.lock.RLock()
	defer strg.lock.RUnlock()
	if clt, ok := strg.sessions[key]; ok {
		return clt.Session, nil
	}
	return nil, nil
}

// CloseSession removes the session from the session registry
func (strg *inMemoryStorage) CloseSession(client *wwr.Client) {
	strg.lock.Lock()
	defer strg.lock.Unlock()
	delete(strg.sessions, client.Session.Key)
}

// hasSession returns true if the given user already has an ongoing session,
// otherwise returns false
func (strg *inMemoryStorage) HasSession(client *wwr.Client) bool {
	strg.lock.RLock()
	defer strg.lock.RUnlock()
	for _, clt := range strg.sessions {
		if clt == client {
			return true
		}
	}
	return false
}

// State represents the servers state
var State = inMemoryStorage{
	lock:      sync.RWMutex{},
	connected: make(map[*wwr.Client]bool),
	sessions:  make(map[string]*wwr.Client),
}
