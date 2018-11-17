package message

import (
	"fmt"
	"io"

	pld "github.com/qbeon/webwire-go/payload"
)

// WriteMsgReply writes a reply message to the given writer closing it
// eventually
func WriteMsgReply(
	writer io.WriteCloser,
	requestIdentifier [8]byte,
	payloadEncoding pld.Encoding,
	payloadData []byte,
) error {
	// Verify payload data validity in case of UTF16 encoding
	if payloadEncoding == pld.Utf16 && len(payloadData)%2 != 0 {
		initialErr := fmt.Errorf(
			"invalid UTF16 reply payload data length: %d",
			len(payloadData),
		)
		if err := writer.Close(); err != nil {
			return fmt.Errorf("%s: %s", initialErr, err)
		}
		return initialErr
	}

	// Determine message type from payload encoding type
	msgType := msgTypeReplyBinary
	if payloadEncoding == pld.Utf8 {
		msgType = msgTypeReplyUtf8
	} else if payloadEncoding == pld.Utf16 {
		msgType = msgTypeReplyUtf16
	}

	// Write message type flag
	if _, err := writer.Write(msgType); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	// Write request identifier
	if _, err := writer.Write(requestIdentifier[:]); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf("%s: %s", err, closeErr)
		}
		return err
	}

	// Write header padding byte if the payload requires proper alignment
	if payloadEncoding == pld.Utf16 {
		if _, err := writer.Write(msgPayloadOffset); err != nil {
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

	if err := writer.Close(); err != nil {
		return fmt.Errorf("%s: %s", err, err)
	}
	return nil
}
