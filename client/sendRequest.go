package client

import (
	"fmt"
	"time"

	webwire "github.com/qbeon/webwire-go"
)

func (clt *Client) sendRequest(
	messageType byte,
	name string,
	payload webwire.Payload,
	timeout time.Duration,
) (webwire.Payload, error) {
	// Require either a name or a payload or both
	if len(name) < 1 && len(payload.Data) < 1 {
		return webwire.Payload{}, webwire.NewProtocolErr(
			fmt.Errorf("Invalid request, request message requires " +
				"either a name, a payload or both but is missing both",
			),
		)
	}

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
