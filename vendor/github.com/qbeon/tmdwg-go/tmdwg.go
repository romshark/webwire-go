package tmdwg

import (
	"fmt"
	"sync"
	"time"
)

// TimedWaitGroup represents a WaitGroup
type TimedWaitGroup struct {
	lock      sync.RWMutex
	err       error
	target    int
	current   int
	completed chan struct{}
	timedout  chan struct{}
}

// NewTimedWaitGroup constructs a new timed wait group and starts the
// timeout counter. Target and timeout arguments must both be bigger 0
func NewTimedWaitGroup(target int, timeout time.Duration) *TimedWaitGroup {
	if target < 1 {
		panic(fmt.Errorf("Invalid progress target: %d", target))
	}

	newTask := &TimedWaitGroup{
		lock:      sync.RWMutex{},
		err:       nil,
		target:    target,
		current:   0,
		completed: make(chan struct{}, 1),
		timedout:  make(chan struct{}, 1),
	}

	time.AfterFunc(timeout, func() {
		newTask.lock.Lock()
		newTask.err = fmt.Errorf(
			"Timed out after %s at the progress: %d of %d",
			timeout,
			newTask.current,
			target,
		)
		newTask.lock.Unlock()
		close(newTask.timedout)
	})

	return newTask
}

// IsCompleted returns true if the wait group already completed,
// otherwise returns false
func (p *TimedWaitGroup) IsCompleted() bool {
	select {
	case <-p.completed:
		return true
	default:
	}
	return false
}

// CurrentProgress returns the current progress
func (p *TimedWaitGroup) CurrentProgress() int {
	p.lock.RLock()
	currentProgress := p.current
	p.lock.RUnlock()
	return currentProgress
}

// Progress progresses the wait group by the given delta.
// Returns the current progress after the update
func (p *TimedWaitGroup) Progress(delta int) int {
	p.lock.Lock()
	if p.IsCompleted() {
		currentProgress := p.current
		p.lock.Unlock()
		return currentProgress
	}

	currentProgress := p.current + delta
	p.current = currentProgress
	if p.current >= p.target {
		close(p.completed)
	}
	p.lock.Unlock()
	return currentProgress
}

// Wait blocks until the wait group is either completed or timed out.
// Returns an error if the wait group timed out
func (p *TimedWaitGroup) Wait() error {
	if p.IsCompleted() {
		p.lock.RLock()
		err := p.err
		p.lock.RUnlock()
		return err
	}

	select {
	case <-p.timedout:
		p.lock.RLock()
		err := p.err
		p.lock.RUnlock()
		return err
	case <-p.completed:
	}

	return nil
}
