package msgbuf

import (
	"errors"
	"io"
	"io/ioutil"
)

// MessageBuffer represents a message buffer
type MessageBuffer struct {
	buf     []byte
	len     int
	onClose func()
}

// Bytes returns a full-length slice of the buffer
func (mb *MessageBuffer) Bytes() []byte {
	return mb.buf
}

// Close resets the message buffer and puts it back into the original pool
func (mb *MessageBuffer) Close() {
	if mb.len < 1 {
		return
	}

	// Reset only the used part of the buffer
	filledSlice := mb.buf[:mb.len]
	for i := range filledSlice {
		filledSlice[i] = 0
	}
	mb.len = 0

	// Call the closure callback
	mb.onClose()
}

// Read reads from the given reader until EOF or error
func (mb *MessageBuffer) Read(reader io.Reader) error {
	cursor := 0
	for {
		if cursor >= len(mb.buf) {
			// Expect EOF on full buffer
			if _, err := reader.Read(mb.buf[:]); err != io.EOF {
				// Overflow! Discard the message that's bigger than the buffer
				mb.Close()
				io.Copy(ioutil.Discard, reader)
				return errors.New("message buffer overflow")
			}

			// Successfully read out the reader
			mb.len = cursor
			return nil
		}

		readBytes, err := reader.Read(mb.buf[cursor:])
		cursor += readBytes

		if readBytes < 0 {
			panic("negative read len")
		}
		if err != nil {
			if err == io.EOF {
				mb.len = cursor
				return nil
			}

			mb.Close()
			return err
		}
	}
}

// Data returns a slice of the usable part of the buffer
func (mb *MessageBuffer) Data() []byte {
	return mb.buf[:mb.len]
}
