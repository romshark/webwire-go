package client

import (
	"context"
	"sync"

	webwire "github.com/qbeon/webwire-go"
)

// dam represents a "goroutine dam" that accumulates goroutines blocking them
// until it's flushed
type dam struct {
	lock    sync.RWMutex
	trigger chan error
	err     error
}

// newDam constructs a new dam instance
func newDam() *dam {
	return &dam{
		lock:    sync.RWMutex{},
		trigger: make(chan error),
		err:     nil,
	}
}

// await blocks the calling goroutine until the dam is flushed
func (dam *dam) await(
	ctx context.Context,
	ctxHasDeadline bool,
) error {
	dam.lock.RLock()
	trigger := dam.trigger
	dam.lock.RUnlock()
	select {
	case <-ctx.Done():
		// Return context error if the context initially had a deadline
		if ctxHasDeadline {
			return ctx.Err()
		}
		// Or return a default timeout if the deadline was set automatically
		return webwire.TimeoutErr{}
	case <-trigger:
		dam.lock.RLock()
		err := dam.err
		dam.lock.RUnlock()
		return err
	}
}

// flush flushes the dam freeing all accumulated goroutines
func (dam *dam) flush(err error) {
	dam.lock.Lock()
	close(dam.trigger)
	dam.err = err
	dam.trigger = make(chan error)
	dam.lock.Unlock()
}
