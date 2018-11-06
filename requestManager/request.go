package requestmanager

import (
	"context"
	"errors"
	"time"

	webwire "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
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

	// timeout represents the configured timeout duration of this request
	timeout time.Duration

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
	// Start timeout timer
	timeoutTimer := time.NewTimer(req.timeout)

	// Block until either deadline exceeded, canceled,
	// timed out or reply received
	select {
	case <-ctx.Done():
		timeoutTimer.Stop()
		req.manager.deregister(req.identifier)
		return nil, webwire.TranslateContextError(ctx.Err())

	case <-timeoutTimer.C:
		timeoutTimer.Stop()
		req.manager.deregister(req.identifier)
		return nil, webwire.NewTimeoutErr(errors.New("timed out"))

	case rp := <-req.reply:
		timeoutTimer.Stop()
		if rp.Error != nil {
			return nil, rp.Error
		}
		return &reply{msg: rp.ReplyMsg}, nil
	}
}
