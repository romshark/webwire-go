package test

import (
	"fmt"
	"sync"
	"time"
)

// pending represents a timed asynchronous task.
type pending struct {
	target  uint32
	current uint32
	timeout time.Duration
	done    chan bool
	barrier chan bool
	timer   *time.Timer
	err     error
	lock    sync.Mutex
}

// newPending returns a new pending asynchronous task.
// Target defines the total required progress.
// Timeout defines the deadline timer duration.
// If start is true the timer is started immediately
func newPending(target uint32, timeout time.Duration, start bool) *pending {
	if target < 1 {
		panic(fmt.Errorf("pending.target cannot be zero"))
	}
	pen := &pending{
		target,
		0,
		timeout,
		make(chan bool, 1),
		make(chan bool, 1),
		nil,
		nil,
		sync.Mutex{},
	}
	if start {
		pen.Start()
	}
	return pen
}

// Start starts the timer. Does nothing if the timer is already running
func (pen *pending) Start() {
	pen.lock.Lock()
	defer pen.lock.Unlock()

	if pen.timer != nil {
		return
	}

	pen.timer = time.NewTimer(pen.timeout)
	go func() {
	LOOP:
		for {
			select {
			case <-pen.timer.C:
				pen.lock.Lock()
				pen.err = fmt.Errorf("Pending task timed out")
				pen.lock.Unlock()
				break LOOP
			case <-pen.done:
				pen.lock.Lock()
				pen.current++
				if pen.current < pen.target {
					pen.lock.Unlock()
					continue
				}
				// Success
				pen.lock.Unlock()
				break LOOP
			}
		}
		pen.lock.Lock()
		pen.done = nil
		close(pen.barrier)
		pen.lock.Unlock()
	}()
}

// Done progresses the task by 1
func (pen *pending) Done() {
	pen.done <- true
}

// Wait blocks until the task is either accomplished or timed out.
// Returns an error if the task timed out
func (pen *pending) Wait() error {
	pen.lock.Lock()
	if pen.done == nil {
		return pen.err
	}
	pen.lock.Unlock()

	<-pen.barrier
	pen.lock.Lock()
	defer pen.lock.Unlock()
	return pen.err
}
