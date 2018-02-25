package test

import (
	"fmt"
	"sync"
	"time"
)

// Pending represents a timed asynchronous task.
type Pending struct {
	target  uint32
	current uint32
	timeout time.Duration
	done    chan bool
	barrier chan bool
	timer   *time.Timer
	result  bool
	lock    sync.Mutex
}

// NewPending returns a new pending asynchronous task.
// Target defines the total required progress.
// Timeout defines the deadline timer duration.
// If start is true the timer is started immediately
func NewPending(target uint32, timeout time.Duration, start bool) *Pending {
	if target < 1 {
		panic(fmt.Errorf("Pending.target cannot be zero"))
	}
	pen := &Pending{
		target,
		0,
		timeout,
		make(chan bool, 1),
		make(chan bool, 1),
		nil,
		false,
		sync.Mutex{},
	}
	if start {
		pen.Start()
	}
	return pen
}

// Start starts the timer. Does nothing if the timer is already running
func (pen *Pending) Start() {
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
				break LOOP
			case <-pen.done:
				pen.lock.Lock()
				pen.current++
				if pen.current < pen.target {
					pen.lock.Unlock()
					continue
				}
				// Success
				pen.result = true
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
func (pen *Pending) Done() {
	pen.done <- true
}

// Wait blocks until the task is either accomplished or timed out.
// Returns an error if the task timed out
func (pen *Pending) Wait() error {
	pen.lock.Lock()
	if pen.done == nil {
		return nil
	}
	pen.lock.Unlock()

	<-pen.barrier
	pen.lock.Lock()
	defer pen.lock.Unlock()
	if pen.result {
		return nil
	}
	return fmt.Errorf("Pending task timed out")
}
