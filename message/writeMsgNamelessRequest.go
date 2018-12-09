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
	identifier []byte,
	binaryPayload []byte,
) error {
	msgType := msgTypeRequestCloseSession
	if reqType == MsgRequestRestoreSession {
		msgType = msgTypeRequestRestoreSession
	} else if reqType != MsgDoCloseSession {
		panic(fmt.Errorf("unexpected nameless request type: %d", reqType))
	}

	// Write message type flag
	if _, err := writer.Write(msgType); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	// Write request identifier
	if _, err := writer.Write(identifier); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	// Write payload
	if len(binaryPayload) > 0 {
		if _, err := writer.Write(binaryPayload); err != nil {
			if closeErr := writer.Close(); closeErr != nil {
				return fmt.Errorf("%s: %s", err, closeErr)
			}
			return err
		}
	}

	return writer.Close()
}
