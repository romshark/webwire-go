package client

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	webwire "github.com/qbeon/webwire-go"
)

func (clt *Client) sendRequest(
	messageType rune,
	payload []byte,
	timeout time.Duration,
) ([]byte, *webwire.Error) {
	if atomic.LoadInt32(&clt.isConnected) < 1 {
		return nil, &webwire.Error{
			Code:    "DISCONNECTED",
			Message: "Trying to send a request on a disconnected socket",
		}
	}

	request := clt.requestManager.Create(timeout)
	reqIdentifier := request.Identifier()

	var msg bytes.Buffer
	msg.WriteRune(messageType)
	msg.Write(reqIdentifier[:])
	msg.Write(payload)

	// Send request
	clt.connLock.Lock()
	err := clt.conn.WriteMessage(websocket.TextMessage, msg.Bytes())
	clt.connLock.Unlock()
	if err != nil {
		// TODO: return typed error TransmissionFailure
		return nil, &webwire.Error{
			Message: fmt.Sprintf("Couldn't send message: %s", err),
		}
	}

	// Block until request either times out or a response is received
	return request.AwaitReply()
}
