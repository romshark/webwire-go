package client

import (
	"context"
	"fmt"

	webwire "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/wwrerr"
)

func (clt *client) sendRequest(
	ctx context.Context,
	messageType byte,
	name []byte,
	payload webwire.Payload,
) (webwire.Reply, error) {
	// Require either a name or a payload or both
	if len(name) < 1 && len(payload.Data) < 1 {
		return nil, wwrerr.ProtocolErr{
			Cause: fmt.Errorf("Invalid request, request message requires " +
				"either a name, a payload or both but is missing both",
			),
		}
	}

	// Register a new request
	request := clt.requestManager.Create()

	writer, err := clt.conn.GetWriter()
	if err != nil {
		return nil, err
	}

	// Compose a message and register it
	if err := message.WriteMsgRequest(
		writer,
		request.IdentifierBytes,
		name,
		payload.Encoding,
		payload.Data,
		true,
	); err != nil {
		clt.requestManager.Fail(request.Identifier, err)
		return nil, wwrerr.TransmissionErr{Cause: err}
	}

	clt.heartbeat.reset()

	// Block until request either times out or a response is received
	return request.AwaitReply(ctx)
}
