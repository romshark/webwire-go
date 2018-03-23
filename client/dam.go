package client

import (
	"sync"
	"time"

	wwr "github.com/qbeon/webwire-go"
)

// dam represents a "goroutine dam" that accumulates goroutines blocking them until it's flushed
type dam struct {
	lock    sync.RWMutex
	barrier chan error
}

// newDam constructs a new dam instance
func newDam() *dam {
	return &dam{
		lock:    sync.RWMutex{},
		barrier: make(chan error),
	}
}

// await blocks the calling goroutine until the dam is flushed
func (dam *dam) await(timeout time.Duration) error {
	dam.lock.RLock()
	defer dam.lock.RUnlock()
	if timeout > 0 {
		select {
		case err := <-dam.barrier:
			return err
		case <-time.After(timeout):
			return wwr.ReqTimeoutErr{}
		}
	} else {
		return <-dam.barrier
	}
}

// flush flushes the dam freeing all accumulated goroutines
func (dam *dam) flush(err error) {
	close(dam.barrier)

	// Reset barrier
	dam.lock.Lock()
	dam.barrier = make(chan error)
	dam.lock.Unlock()
}
