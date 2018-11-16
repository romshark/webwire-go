package message

import (
	"fmt"
	"io"
)

// WriteMsgNamelessRequest writes a nameless (initially without a name)
// request message to the given writer closing it eventually
func WriteMsgNamelessRequest(
	writer io.WriteCloser,
	reqType byte,
	identifier [8]byte,
	binaryPayload []byte,
) error {
	// Write message type flag
	if _, err := writer.Write([]byte{reqType}); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	// Write request identifier
	if _, err := writer.Write(identifier[:]); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	// Write payload
	if _, err := writer.Write(binaryPayload); err != nil {
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
