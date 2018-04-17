package test

import (
	"testing"
	"time"
)

// TestPendingInvalidTarget tests passing invalid zero target to
// pending constructor
func TestPendingInvalidTarget(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("Expected a panic after passing invalid parameters")
		}
	}()

	newPending(0, 1*time.Second, true)
}

// TestPendingWaitAfterDone tests the pending primitive
// to pass through wait when already done
func TestPendingWaitAfterDone(t *testing.T) {
	checkpoint := newPending(1, 1*time.Second, true)
	checkpoint.Done()

	if err := checkpoint.Wait(); err != nil {
		t.Fatalf("Expected done pending to have no error, got: %s", err)
	}
}

// TestPendingWaitAfterTimeout tests the pending primitive
// to pass through wait when already done
func TestPendingWaitAfterTimeout(t *testing.T) {
	checkpoint := newPending(1, 10*time.Millisecond, true)

	// Wait for the timeout
	time.Sleep(100 * time.Millisecond)

	if err := checkpoint.Wait(); err == nil {
		t.Fatal("Expected pending to have timed out before")
	}
}

// TestPendingWaitTimeout tests the pending primitive's Wait method
// expecting it to return a timeout error
func TestPendingWaitTimeout(t *testing.T) {
	checkpoint := newPending(1, 1*time.Millisecond, true)

	// Wait for 1 millisecond, then timeout and return err
	if err := checkpoint.Wait(); err == nil {
		t.Fatal("Expected pending to return an error")
	}
}

// TestPendingWaitMultipleTimeout tests timeout of progress of 2
func TestPendingWaitMultipleTimeout(t *testing.T) {
	checkpoint := newPending(2, 1*time.Millisecond, true)
	checkpoint.Done()

	// Wait for 1 millisecond, then timeout and return err
	if err := checkpoint.Wait(); err == nil {
		t.Fatal("Expected pending to return an error")
	}
}
