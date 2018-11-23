package message

import (
	"errors"
	"fmt"
	"io"

	pld "github.com/qbeon/webwire-go/payload"
)

// WriteMsgRequest writes a named request message to the given writer
// closing it eventually
func WriteMsgRequest(
	writer io.WriteCloser,
	identifier []byte,
	name []byte,
	payloadEncoding pld.Encoding,
	payloadData []byte,
	safeMode bool,
) error {
	// Require either a name, or a payload or both, but don't allow none
	if len(name) < 1 && len(payloadData) < 1 {
		initialErr := errors.New(
			"request message requires either a name, or a payload, or both",
		)
		if err := writer.Close(); err != nil {
			return fmt.Errorf("%s: %s", initialErr, err)
		}
		return initialErr
	}

	// Cap name length at 255 bytes
	if len(name) > 255 {
		initialErr := fmt.Errorf(
			"unsupported request message name length: %d",
			len(name),
		)
		if err := writer.Close(); err != nil {
			return fmt.Errorf("%s: %s", initialErr, err)
		}
		return initialErr
	}

	// Verify payload data validity in case of UTF16 encoding
	if payloadEncoding == pld.Utf16 && len(payloadData)%2 != 0 {
		initialErr := fmt.Errorf(
			"invalid UTF16 request payload data length: %d",
			len(payloadData),
		)
		if err := writer.Close(); err != nil {
			return fmt.Errorf("%s: %s", initialErr, err)
		}
		return initialErr
	}

	// Validate name
	if safeMode {
		for i := 0; i < len(name); i++ {
			char := name[i]
			if char < 32 || char > 126 {
				initialErr := fmt.Errorf(
					"unsupported character in request name: %s",
					string(char),
				)
				if err := writer.Close(); err != nil {
					return fmt.Errorf("%s: %s", initialErr, err)
				}
				return initialErr
			}
		}
	}

	// Determine message type from payload encoding type
	msgType := msgTypeRequestBinary
	if payloadEncoding == pld.Utf8 {
		msgType = msgTypeRequestUtf8
	} else if payloadEncoding == pld.Utf16 {
		msgType = msgTypeRequestUtf16
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

	// Write name length flag
	if _, err := writer.Write(msgNameLenBytes[len(name)]); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}
	// Write name
	if _, err := writer.Write(name); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	// Write header padding byte if the payload requires proper alignment
	if payloadEncoding == pld.Utf16 && len(name)%2 != 0 {
		if _, err := writer.Write(msgHeaderPadding); err != nil {
			if closeErr := writer.Close(); closeErr != nil {
				return fmt.Errorf("%s: %s", err, closeErr)
			}
			return err
		}
	}

	// Write payload
	if _, err := writer.Write(payloadData); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	return writer.Close()
}
