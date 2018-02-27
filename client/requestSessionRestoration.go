package client

import (
	"encoding/json"
	"fmt"

	webwire "github.com/qbeon/webwire-go"
)

// requestSessionRestoration sends a session restoration request
// and decodes the session object from the received replied.
// Expects the client to be connected beforehand
func (clt *Client) requestSessionRestoration(sessionKey []byte) (*webwire.Session, error) {
	reply, err := clt.sendRequest(
		webwire.MsgRestoreSession,
		sessionKey,
		clt.defaultTimeout,
	)
	if err != nil {
		// TODO: check for error types
		return nil, fmt.Errorf("Session restoration request failed: %s", err)
	}

	var session webwire.Session
	if err := json.Unmarshal(reply, &session); err != nil {
		return nil, fmt.Errorf(
			"Couldn't unmarshal restored session from reply('%s'): %s",
			reply,
			err,
		)
	}

	return &session, nil
}
