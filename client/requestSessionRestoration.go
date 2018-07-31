package client

import (
	"context"
	"encoding/json"
	"fmt"

	webwire "github.com/qbeon/webwire-go"
	msg "github.com/qbeon/webwire-go/message"
	pld "github.com/qbeon/webwire-go/payload"
)

// requestSessionRestoration sends a session restoration request
// and decodes the session object from the received reply.
// Expects the client to be connected beforehand
func (clt *Client) requestSessionRestoration(sessionKey []byte) (
	*webwire.Session,
	error,
) {
	reply, err := clt.sendNamelessRequest(
		context.Background(),
		msg.MsgRestoreSession,
		pld.Payload{
			Encoding: webwire.EncodingBinary,
			Data:     sessionKey,
		},
		clt.defaultReqTimeout,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON encoded session object
	var encodedSessionObj webwire.JSONEncodedSession
	if err := json.Unmarshal(reply.Data(), &encodedSessionObj); err != nil {
		return nil, fmt.Errorf(
			"Couldn't unmarshal restored session from reply('%s'): %s",
			string(reply.Data()),
			err,
		)
	}

	// Parse session info object
	var decodedInfo webwire.SessionInfo
	if encodedSessionObj.Info != nil {
		decodedInfo = clt.sessionInfoParser(encodedSessionObj.Info)
	}

	return &webwire.Session{
		Key:      encodedSessionObj.Key,
		Creation: encodedSessionObj.Creation,
		Info:     decodedInfo,
	}, nil
}
