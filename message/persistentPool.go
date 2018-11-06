package message

import "sync"

// PersistentPool represents a thread-safe messageBuffer pool
type PersistentPool struct {
	bufferSize uint32
	lock       sync.Mutex
	avail      map[*Message]struct{}
}

// create creates a new pooled message instance
func (pl *PersistentPool) create() *Message {
	newMessage := NewMessage(pl.bufferSize)
	newMessage.onClose = func() {
		pl.lock.Lock()
		pl.avail[newMessage] = struct{}{}
		pl.lock.Unlock()
	}
	return newMessage
}

// NewPersistentPool initializes a new persistent message buffer pool instance
func NewPersistentPool(bufferSize, prealloc uint32) *PersistentPool {
	newPool := &PersistentPool{
		bufferSize: bufferSize,
		lock:       sync.Mutex{},
		avail:      make(map[*Message]struct{}, prealloc),
	}
	for i := uint32(0); i < prealloc; i++ {
		newPool.avail[newPool.create()] = struct{}{}
	}
	return newPool
}

// Get implements the Pool interface
func (pl *PersistentPool) Get() *Message {
	pl.lock.Lock()
	if len(pl.avail) > 0 {
		for msg := range pl.avail {
			pl.lock.Unlock()
			return msg
		}
	}
	pl.lock.Unlock()

	// Allocate a new message object
	return pl.create()
}

// Purge purges all currently unused messages
func (pl *PersistentPool) Purge() {
	pl.lock.Lock()
	pl.avail = make(map[*Message]struct{})
	pl.lock.Unlock()
}
