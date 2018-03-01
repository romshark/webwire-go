package client

import (
	"encoding/json"

	webwire "github.com/qbeon/webwire-go"
)

func (clt *Client) handleSessionCreated(sessionKey []byte) {
	// Set new session
	var session webwire.Session

	if err := json.Unmarshal(sessionKey, &session); err != nil {
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

func (clt *Client) handleFailure(reqID [8]byte, payload []byte) {
	// Decode error
	var replyErr webwire.Error
	if err := json.Unmarshal(payload, &replyErr); err != nil {
		clt.errorLog.Printf("Failed unmarshalling error reply: %s", err)
	}

	// Fail request
	clt.requestManager.Fail(reqID, replyErr)
}

func (clt *Client) handleReply(reqID [8]byte, payload webwire.Payload) {
	clt.requestManager.Fulfill(reqID, payload)
}

func (clt *Client) handleMessage(message []byte) error {
	if len(message) < 1 {
		return nil
	}
	switch message[0:1][0] {
	case webwire.MsgReplyBinary:
		clt.handleReply(
			extractMessageIdentifier(message),
			webwire.Payload{
				Encoding: webwire.EncodingBinary,
				Data:     message[9:],
			},
		)
	case webwire.MsgReplyUtf8:
		clt.handleReply(
			extractMessageIdentifier(message),
			webwire.Payload{
				Encoding: webwire.EncodingUtf8,
				Data:     message[9:],
			},
		)
	case webwire.MsgReplyUtf16:
		clt.handleReply(
			extractMessageIdentifier(message),
			webwire.Payload{
				Encoding: webwire.EncodingUtf16,
				Data:     message[9:],
			},
		)
	case webwire.MsgErrorReply:
		clt.handleFailure(extractMessageIdentifier(message), message[9:])
	case webwire.MsgSignalBinary:
		clt.hooks.OnServerSignal(webwire.Payload{
			Encoding: webwire.EncodingBinary,
			Data:     message[2:],
		})
	case webwire.MsgSignalUtf8:
		clt.hooks.OnServerSignal(webwire.Payload{
			Encoding: webwire.EncodingUtf8,
			Data:     message[2:],
		})
	case webwire.MsgSignalUtf16:
		clt.hooks.OnServerSignal(webwire.Payload{
			Encoding: webwire.EncodingUtf16,
			Data:     message[2:],
		})
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
