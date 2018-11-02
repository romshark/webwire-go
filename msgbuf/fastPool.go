package msgbuf

import (
	"sync"
)

// FastPool represents a thread-safe messageBuffer pool
type FastPool struct {
	maxMsgSize uint32
	lock       sync.Mutex
	available  map[*MessageBuffer]struct{}
	total      uint64
}

// NewFastPool initializes a new message buffer pool instance
func NewFastPool(maxMsgSize, prealloc uint32) *FastPool {
	fp := &FastPool{
		maxMsgSize: maxMsgSize,
		lock:       sync.Mutex{},
		available:  make(map[*MessageBuffer]struct{}, prealloc),
		total:      0,
	}

	for i := uint32(0); i < prealloc; i++ {
		// Allocate a new message buffer
		newMessageBuffer := &MessageBuffer{
			buf: make([]byte, fp.maxMsgSize),
		}
		fp.available[newMessageBuffer] = struct{}{}
	}

	return fp
}

// Get implements the Pool interface
func (fp *FastPool) Get() *MessageBuffer {
	fp.lock.Lock()
	if len(fp.available) < 1 {
		// Allocate a new message buffer
		newMessageBuffer := &MessageBuffer{
			buf: make([]byte, fp.maxMsgSize),
		}
		newMessageBuffer.onClose = func() {
			// Put the message buffer back into the pool
			fp.lock.Lock()
			fp.available[newMessageBuffer] = struct{}{}
			fp.lock.Unlock()
		}
		fp.total++
		fp.lock.Unlock()
		return newMessageBuffer
	}

	// Take one of the available message buffers from the pool
	var taken *MessageBuffer
	for buf := range fp.available {
		taken = buf
		break
	}
	// Remove it from the register of available buffers
	delete(fp.available, taken)
	fp.lock.Unlock()

	return taken
}
