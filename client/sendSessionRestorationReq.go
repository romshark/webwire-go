package client

import (
	"fmt"

	webwire "github.com/qbeon/webwire-go"
)

// sendSessionRestorationReq sends a session restoration request.
// Expects the client to be connected beforehand
func (clt *Client) sendSessionRestorationReq(sessionKey []byte) error {
	if _, err := clt.sendRequest(
		webwire.MsgRestoreSession,
		sessionKey,
		clt.defaultTimeout,
	); err != nil {
		// TODO: check for error types
		return fmt.Errorf("Session restoration request failed: %s", err)
	}
	return nil
}
