package message

import (
	"fmt"
	"io"
)

// WriteMsgEmptyRequest writes a request message consisting only of the type and
// identifier to the given writer closing it eventually
func WriteMsgEmptyRequest(
	writer io.WriteCloser,
	msgType byte,
	id [8]byte,
) error {
	// Write message type flag
	if _, err := writer.Write([]byte{msgType}); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	// Write request identifier
	if _, err := writer.Write(id[:]); err != nil {
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
