package requestmanager

import (
	"context"

	webwire "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/wwrerr"
)

// TODO: The request identifier should remain a uint64 until it's converted into
// the byte array for transmission, this would slightly increase performance

// RequestIdentifier represents the identifier of a request.
type RequestIdentifier = [8]byte

// genericReply is used by the request manager to represent the results of a
// request that either failed or succeeded
type genericReply struct {
	ReplyMsg *message.Message
	Error    error
}

// Request represents a request created and tracked by the request manager
type Request struct {
	// manager references the RequestManager instance managing this request
	manager *RequestManager

	// identifier represents the unique identifier of this request
	identifier RequestIdentifier

	// reply represents a channel for asynchronous reply handling
	reply chan genericReply
}

// Identifier returns the assigned request identifier
func (req *Request) Identifier() RequestIdentifier {
	return req.identifier
}

// AwaitReply blocks the calling goroutine
// until either the reply is fulfilled or failed, the request timed out
// a user-defined deadline was exceeded or the request was prematurely canceled.
// The timer is started when AwaitReply is called.
func (req *Request) AwaitReply(ctx context.Context) (webwire.Reply, error) {
	// Block until either context canceled (including timeout) or reply received
	select {
	case <-ctx.Done():
		req.manager.deregister(req.identifier)
		return nil, wwrerr.TranslateContextError(ctx.Err())

	case rp := <-req.reply:
		if rp.Error != nil {
			return nil, rp.Error
		}
		return &reply{msg: rp.ReplyMsg}, nil
	}
}
