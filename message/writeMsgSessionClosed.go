package message

import (
	"fmt"
	"io"
)

// WriteMsgSessionClosed writes a session closure notification message to the
// given writer closing it eventually
func WriteMsgSessionClosed(writer io.WriteCloser) error {
	// Write message type flag
	if _, err := writer.Write(msgTypeSessionClosed); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	return writer.Close()
}
