package memchan

import (
	"errors"
)

// Buffer represents a reactive outbound buffer implementation
type Buffer struct {
	buf     []byte
	len     int
	onFlush func([]byte) error
}

// NewBuffer allocates a new buffer
func NewBuffer(buf []byte, onFlush func([]byte) error) Buffer {
	if len(buf) < 1 {
		panic("empty buffer")
	}
	if onFlush == nil {
		onFlush = func([]byte) error { return nil }
	}
	return Buffer{
		buf:     buf,
		len:     0,
		onFlush: onFlush,
	}
}

// reset clears the buffer
func (buf *Buffer) reset() {
	buf.len = 0
}

// Write writes a portion of data to the buffer
func (buf *Buffer) Write(p []byte) (int, error) {
	if len(p) > len(buf.buf)-buf.len {
		// Buffer overflow
		buf.reset()
		return 0, errors.New("buffer overflow")
	}
	copy(buf.buf[buf.len:], p)
	buf.len += len(p)
	return len(p), nil
}

// Close flushes the buffer to the reader
func (buf *Buffer) Close() (err error) {
	err = buf.onFlush(buf.buf[:buf.len])
	buf.reset()
	return err
}
