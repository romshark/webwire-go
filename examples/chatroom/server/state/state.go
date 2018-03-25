package state

import (
	"sync"

	wwr "github.com/qbeon/webwire-go"
)

type inMemoryStorage struct {
	lock      sync.RWMutex
	connected map[*wwr.Client]bool
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

// State represents the servers state
var State = inMemoryStorage{
	lock:      sync.RWMutex{},
	connected: make(map[*wwr.Client]bool),
}
