package message

import (
	"fmt"
	"io"
)

// WriteMsgErrorReply writes an error reply message to the given writer
// closing it eventually
func WriteMsgErrorReply(
	writer io.WriteCloser,
	requestIdent []byte,
	code,
	message []byte,
	safeMode bool,
) error {
	if safeMode {
		// Validate input
		if len(code) < 1 {
			initialErr := fmt.Errorf(
				"missing error code while creating a new error reply message",
			)
			if err := writer.Close(); err != nil {
				return fmt.Errorf("%s: %s", initialErr, err)
			}
			return initialErr
		} else if len(code) > 255 {
			initialErr := fmt.Errorf(
				"invalid error code while creating a new error reply message,"+
					"too long (%d)",
				len(code),
			)
			if err := writer.Close(); err != nil {
				return fmt.Errorf("%s: %s", initialErr, err)
			}
			return initialErr
		}
		// Determine total message length
		// messageSize := 10 + len(code) + len(message)
		// if len(buf) < messageSize {
		//	if closeErr := writer.Close(); closeErr != nil {
		// 	return fmt.Errorf("%s: %s", err, closeErr)
		// }
		// 	return 0, errors.New(
		// 		"message buffer too small to fit an error reply message",
		// 	)
		// }
		for _, char := range code {
			if char < 32 || char > 126 {
				initialErr := fmt.Errorf(
					"unsupported character in reply error - error code: %s",
					string(char),
				)
				if closeErr := writer.Close(); closeErr != nil {
					return fmt.Errorf("%s: %s", initialErr, closeErr)
				}
				return initialErr
			}
		}
	}

	// Write message type flag
	if _, err := writer.Write(msgTypeReplyError); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	// Write request identifier
	if _, err := writer.Write(requestIdent); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	// Write code length flag
	if _, err := writer.Write(
		msgNameLenBytes[len(code) : len(code)+1],
	); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	// Write error code
	if _, err := writer.Write(code); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	// Write error message
	if _, err := writer.Write(message); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	return writer.Close()
}
