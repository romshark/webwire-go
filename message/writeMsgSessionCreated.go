package message

import (
	"fmt"
	"io"
)

// WriteMsgNotifySessionCreated writes a session creation notification message to
// the given writer closing it eventually
func WriteMsgNotifySessionCreated(
	writer io.WriteCloser,
	sessionInfo []byte,
) error {
	// Write message type flag
	if _, err := writer.Write(msgTypeSessionCreated); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	// Write the session info payload
	if _, err := writer.Write(sessionInfo); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	return writer.Close()
}
