package client

import (
	"time"

	webwire "github.com/qbeon/webwire-go"
)

func (clt *Client) sendNamelessRequest(
	messageType byte,
	payload webwire.Payload,
	timeout time.Duration,
) (webwire.Payload, error) {
	request := clt.requestManager.Create(timeout)
	reqIdentifier := request.Identifier()

	msg := webwire.NewNamelessRequestMessage(messageType, reqIdentifier, payload.Data)

	// Send request
	if err := clt.conn.Write(msg); err != nil {
		return webwire.Payload{}, webwire.NewReqTransErr(err)
	}

	// Block until request either times out or a response is received
	return request.AwaitReply()
}
