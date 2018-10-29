package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	wwr "github.com/qbeon/webwire-go"
)

// dam represents a "goroutine dam" that accumulates goroutines blocking them
// until it's flushed
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
func (dam *dam) await(ctx context.Context, timeout time.Duration) error {
	dam.lock.RLock()
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		select {
		case <-ctx.Done():
			dam.lock.RUnlock()
			timer.Stop()
			return wwr.TranslateContextError(ctx.Err())
		case err := <-dam.barrier:
			dam.lock.RUnlock()
			timer.Stop()
			return err
		case <-timer.C:
			dam.lock.RUnlock()
			timer.Stop()
			return wwr.NewTimeoutErr(fmt.Errorf("timed out"))
		}
	} else {
		dam.lock.RUnlock()
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
