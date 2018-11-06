package message

import (
	"fmt"
	"io"
)

// WriteMsgHeartbeat writes a session closure notification message to the
// given writer closing it eventually
func WriteMsgHeartbeat(writer io.WriteCloser) error {
	// Write message type flag
	if _, err := writer.Write([]byte{MsgHeartbeat}); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}
