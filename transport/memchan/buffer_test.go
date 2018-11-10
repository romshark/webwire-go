package memchan

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBufferWrite tests Buffer.Write
func TestBufferWrite(t *testing.T) {
	buf := make([]byte, 6)
	buffer := NewBuffer(buf, nil)
	bytesWritten, err := buffer.Write([]byte{1})
	require.NoError(t, err)
	require.Equal(t, 1, bytesWritten)

	bytesWritten, err = buffer.Write([]byte{2, 3, 4})
	require.NoError(t, err)
	require.Equal(t, 3, bytesWritten)

	require.Equal(t, []byte{1, 2, 3, 4, 0, 0}, buf)

	// Try buffer overflow
	bytesWritten, err = buffer.Write([]byte{1, 2, 3})
	require.Equal(t, 0, bytesWritten)
	require.Error(t, err)
	require.Equal(t, []byte{1, 2, 3, 4, 0, 0}, buf)
}

// TestBufferWriteOverflow tests Buffer.Write overflowing the buffer
func TestBufferWriteOverflow(t *testing.T) {
	buf := make([]byte, 4)
	buffer := NewBuffer(buf, nil)

	bytesWritten, err := buffer.Write([]byte{1, 2, 3, 4, 5})
	require.Equal(t, 0, bytesWritten)
	require.Error(t, err)

	require.Equal(t, []byte{0, 0, 0, 0}, buf)
}

// TestBufferClose tests Buffer.Close
func TestBufferClose(t *testing.T) {
	flushed := false

	buf := make([]byte, 3)
	buffer := NewBuffer(buf, func(data []byte) error {
		flushed = true
		require.Equal(t, []byte{1, 1}, data)
		return nil
	})

	bytesWritten, err := buffer.Write([]byte{1, 1})
	require.Equal(t, 2, bytesWritten)
	require.NoError(t, err)

	require.NoError(t, buffer.Close())
	require.Equal(t, true, flushed)
}

// TestBufferCloseError tests Buffer.Close with error
func TestBufferCloseError(t *testing.T) {
	buf := make([]byte, 3)
	buffer := NewBuffer(buf, func(data []byte) error {
		return errors.New("test error")
	})

	require.Error(t, buffer.Close())
}
