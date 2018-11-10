package message

import "sync"

// SyncPool represents a thread-safe messageBuffer pool
type SyncPool struct {
	bufferSize uint32
	pool       *sync.Pool
}

// NewSyncPool initializes a new sync.Pool based message buffer pool instance
func NewSyncPool(bufferSize, prealloc uint32) *SyncPool {
	pool := &sync.Pool{}
	pool.New = func() interface{} {
		msg := NewMessage(bufferSize)
		msg.onClose = func() {
			pool.Put(msg)
		}
		return msg
	}
	return &SyncPool{
		bufferSize: bufferSize,
		pool:       pool,
	}
}

// Get implements the Pool interface
func (mbp *SyncPool) Get() *Message {
	return mbp.pool.Get().(*Message)
}
