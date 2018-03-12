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
	reply, err := clt.sendNamelessRequest(
		webwire.MsgRestoreSession,
		webwire.Payload{
			Encoding: webwire.EncodingBinary,
			Data:     sessionKey,
		},
		clt.defaultTimeout,
	)
	if err != nil {
		// TODO: check for error types
		fmt.Println("ERR", err)
		return nil, fmt.Errorf("Session restoration request failed: %s", err)
	}

	var session webwire.Session
	if err := json.Unmarshal(reply.Data, &session); err != nil {
		return nil, fmt.Errorf(
			"Couldn't unmarshal restored session from reply('%s'): %s",
			string(reply.Data),
			err,
		)
	}

	return &session, nil
}
