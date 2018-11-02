package client

import (
	"fmt"
)

// verifyProtocolVersion returns true if the given version of the webwire
// protocol is acceptable for this client
func verifyProtocolVersion(major, minor byte) error {
	if major != 2 {
		return fmt.Errorf(
			"unsupported protocol version: %d.%d (supported: 2.x)",
			major,
			minor,
		)
	}
	return nil
}
