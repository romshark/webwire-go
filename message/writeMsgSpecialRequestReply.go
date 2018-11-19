package message

import (
	"fmt"
	"io"
)

// WriteMsgSpecialRequestReply writes a special request reply message to the
// given writer closing it eventually
func WriteMsgSpecialRequestReply(
	writer io.WriteCloser,
	reqType byte,
	reqIdent []byte,
) error {
	msgType := msgTypeReplyInternalError
	switch reqType {
	case MsgInternalError:
	case MsgMaxSessConnsReached:
		msgType = msgTypeReplyMaxSessConnsReached
	case MsgSessionNotFound:
		msgType = msgTypeReplySessionNotFound
	case MsgSessionsDisabled:
		msgType = msgTypeSessionsDisabled
	case MsgReplyShutdown:
		msgType = msgTypeReplyShutdown
	default:
		initialErr := fmt.Errorf(
			"message type (%d) doesn't represent a special reply message",
			reqType,
		)
		if err := writer.Close(); err != nil {
			return fmt.Errorf("%s: %s", initialErr, err)
		}
		return initialErr
	}

	// Write message type flag
	if _, err := writer.Write(msgType); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	// Write request identifier
	if _, err := writer.Write(reqIdent); err != nil {
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
