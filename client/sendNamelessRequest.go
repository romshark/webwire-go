package client

import (
	"time"

	webwire "github.com/qbeon/webwire-go"
	msg "github.com/qbeon/webwire-go/message"
	pld "github.com/qbeon/webwire-go/payload"
)

func (clt *Client) sendNamelessRequest(
	messageType byte,
	payload pld.Payload,
	timeout time.Duration,
) (webwire.Payload, error) {
	request := clt.requestManager.Create(timeout)
	reqIdentifier := request.Identifier()

	msg := msg.NewNamelessRequestMessage(
		messageType,
		reqIdentifier,
		payload.Data,
	)

	// Send request
	if err := clt.conn.Write(msg); err != nil {
		return nil, webwire.NewReqTransErr(err)
	}

	// Block until request either times out or a response is received
	return request.AwaitReply()
}
