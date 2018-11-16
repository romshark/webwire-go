package message

import (
	"errors"
	"io"
	"io/ioutil"
)

// Buffer represents a message buffer
type Buffer struct {
	buf []byte
	len int
}

// Bytes returns a full-length slice of the buffer
func (buf *Buffer) Bytes() []byte {
	return buf.buf
}

// IsEmpty returns true if the buffer is empty, otherwise returns false
func (buf *Buffer) IsEmpty() bool {
	return buf.len < 1
}

// Close resets the message buffer and puts it back into the original pool
func (buf *Buffer) Close() {
	buf.len = 0
}

// Read reads from the given reader until EOF or error
func (buf *Buffer) Read(reader io.Reader) error {
	cursor := 0
	for {
		if cursor >= len(buf.buf) {
			// Expect EOF on full buffer
			_, err := reader.Read(buf.buf)
			if err != io.EOF {
				// Overflow! Discard the message that's bigger than the buffer
				buf.Close()
				io.Copy(ioutil.Discard, reader)
				return errors.New("message buffer overflow")
			}

			// Successfully read out the reader
			buf.len = cursor
			return nil
		}

		readBytes, err := reader.Read(buf.buf[cursor:])
		cursor += readBytes

		if readBytes < 0 {
			panic("negative read len")
		}
		if err != nil {
			if err == io.EOF {
				buf.len = cursor
				return nil
			}

			buf.Close()
			return err
		}
	}
}

// Data returns a slice of the usable part of the buffer
func (buf *Buffer) Data() []byte {
	return buf.buf[:buf.len]
}
