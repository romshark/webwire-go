package msgbuf

import (
	"sync"
)

// SyncPool represents a thread-safe messageBuffer pool
type SyncPool struct {
	maxMsgSize uint32
	pool       *sync.Pool
}

// NewSyncPool initializes a new sync.Pool based message buffer pool instance
func NewSyncPool(maxMsgSize, prealloc uint32) *SyncPool {
	pool := &sync.Pool{}
	pool.New = func() interface{} {
		msgBuf := &MessageBuffer{
			buf: make([]byte, maxMsgSize),
		}
		msgBuf.onClose = func() {
			pool.Put(msgBuf)
		}
		return msgBuf
	}
	return &SyncPool{
		maxMsgSize: maxMsgSize,
		pool:       pool,
	}
}

// Get implements the Pool interface
func (mbp *SyncPool) Get() *MessageBuffer {
	return mbp.pool.Get().(*MessageBuffer)
}
