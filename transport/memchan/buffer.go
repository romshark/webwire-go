package memchan

import (
	"errors"
	"sync"

	wwr "github.com/qbeon/webwire-go"
)

// Buffer represents a reactive outbound buffer implementation
type Buffer struct {
	buf     []byte
	len     int
	onFlush func([]byte) error
	lock    *sync.Mutex
}

// NewBuffer allocates a new buffer
func NewBuffer(
	buf []byte,
	onFlush func([]byte) error,
) Buffer {
	if len(buf) < 1 {
		panic("empty buffer")
	}
	return Buffer{
		buf:     buf,
		len:     0,
		onFlush: onFlush,
		lock:    &sync.Mutex{},
	}
}

// reset clears the buffer
func (buf *Buffer) reset() {
	buf.len = 0
}

// Write writes a portion of data to the buffer
func (buf *Buffer) Write(p []byte) (int, error) {
	buf.lock.Lock()
	if len(p) > len(buf.buf)-buf.len {
		// Buffer overflow
		buf.reset()
		buf.lock.Unlock()
		return 0, wwr.BufferOverflowErr{}
	}
	copy(buf.buf[buf.len:], p)
	buf.len += len(p)
	buf.lock.Unlock()
	return len(p), nil
}

// Close flushes the buffer to the reader
func (buf *Buffer) Close() (err error) {
	buf.lock.Lock()
	if buf.len < 1 {
		return errors.New("no data written")
	}
	err = buf.onFlush(buf.buf[:buf.len])
	buf.reset()
	buf.lock.Unlock()
	return err
}
