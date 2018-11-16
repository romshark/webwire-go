package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	pld "github.com/qbeon/webwire-go/payload"
)

// requestSessionRestoration sends a session restoration request
// and decodes the session object from the received reply.
// Expects the client to be connected beforehand
func (clt *client) requestSessionRestoration(
	ctx context.Context,
	sessionKey []byte,
) (
	*webwire.Session,
	error,
) {
	reply, err := clt.sendNamelessRequest(
		ctx,
		message.MsgRestoreSession,
		pld.Payload{
			Encoding: webwire.EncodingBinary,
			Data:     sessionKey,
		},
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON encoded session object
	var encodedSessionObj webwire.JSONEncodedSession
	if err := json.Unmarshal(
		reply.Payload(),
		&encodedSessionObj,
	); err != nil {
		reply.Close()
		return nil, fmt.Errorf(
			"couldn't unmarshal restored session from reply: %s",
			err,
		)
	}

	reply.Close()

	// Parse session info object
	var decodedInfo webwire.SessionInfo
	if encodedSessionObj.Info != nil {
		decodedInfo = clt.options.SessionInfoParser(encodedSessionObj.Info)
	}

	return &webwire.Session{
		Key:      encodedSessionObj.Key,
		Creation: encodedSessionObj.Creation,
		Info:     decodedInfo,
	}, nil
}
