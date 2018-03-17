package client

import (
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	webwire "github.com/qbeon/webwire-go"
)

func (clt *Client) sendRequest(
	messageType byte,
	name string,
	payload webwire.Payload,
	timeout time.Duration,
) (webwire.Payload, error) {
	if atomic.LoadInt32(&clt.isConnected) < 1 {
		return webwire.Payload{}, webwire.ReqErr{
			Code:    "DISCONNECTED",
			Message: "Trying to send a request on a disconnected socket",
		}
	}

	request := clt.requestManager.Create(timeout)
	reqIdentifier := request.Identifier()

	msg := webwire.NewRequestMessage(reqIdentifier, name, payload)

	// Send request
	clt.connLock.Lock()
	err := clt.conn.WriteMessage(websocket.BinaryMessage, msg)
	clt.connLock.Unlock()
	if err != nil {
		return webwire.Payload{}, webwire.NewReqErrTrans(err)
	}

	// Block until request either times out or a response is received
	return request.AwaitReply()
}
