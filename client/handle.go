package client

import (
	"encoding/json"

	webwire "github.com/qbeon/webwire-go"
)

func (clt *Client) handleSessionCreated(message []byte) {
	// Set new session
	var session webwire.Session

	if err := json.Unmarshal(message, &session); err != nil {
		clt.errorLog.Printf("Failed unmarshalling session object: %s", err)
		return
	}

	clt.session = &session
	clt.hooks.OnSessionCreated(&session)
}

func (clt *Client) handleSessionClosed() {
	// Destroy local session
	clt.session = nil

	clt.hooks.OnSessionClosed()
}

func (clt *Client) handleFailure(message []byte) {
	// Decode error
	var replyErr webwire.Error
	if err := json.Unmarshal(message[33:], &replyErr); err != nil {
		clt.errorLog.Printf("Failed unmarshalling error reply: %s", err)
	}

	// Fail request
	clt.requestManager.Fail(extractMessageIdentifier(message), replyErr)
}

func (clt *Client) handleReply(message []byte) {
	clt.requestManager.Fulfill(extractMessageIdentifier(message), message[33:])
}

func (clt *Client) handleMessage(message []byte) error {
	if len(message) < 1 {
		return nil
	}
	switch message[0:1][0] {
	case webwire.MsgReply:
		clt.handleReply(message)
	case webwire.MsgErrorReply:
		clt.handleFailure(message)
	case webwire.MsgSignal:
		clt.hooks.OnServerSignal(message[1:])
	case webwire.MsgSessionCreated:
		clt.handleSessionCreated(message[1:])
	case webwire.MsgSessionClosed:
		clt.handleSessionClosed()
	default:
		clt.warningLog.Printf(
			"Strange message type received: '%c'\n",
			message[0:1][0],
		)
	}
	return nil
}
