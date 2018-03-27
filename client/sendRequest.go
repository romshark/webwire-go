package client

import (
	"time"

	webwire "github.com/qbeon/webwire-go"
)

func (clt *Client) sendRequest(
	messageType byte,
	name string,
	payload webwire.Payload,
	timeout time.Duration,
) (webwire.Payload, error) {
	request := clt.requestManager.Create(timeout)
	reqIdentifier := request.Identifier()

	msg := webwire.NewRequestMessage(reqIdentifier, name, payload)

	// Send request
	if err := clt.conn.Write(msg); err != nil {
		return webwire.Payload{}, webwire.NewReqTransErr(err)
	}

	// Block until request either times out or a response is received
	return request.AwaitReply()
}
