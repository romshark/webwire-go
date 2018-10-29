package requestmanager

import (
	"context"
	"fmt"
	"time"

	webwire "github.com/qbeon/webwire-go"
)

// TODO: The request identifier should remain a uint64 until it's converted into
// the byte array for transmission, this would slightly increase performance

// RequestIdentifier represents the universally unique, minified
// UUIDv4 identifier of a request.
type RequestIdentifier = [8]byte

// reply is used by the request manager to represent the results
// of a request (both failed and succeeded)
type reply struct {
	Reply webwire.Payload
	Error error
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
	reply chan reply
}

// Identifier returns the assigned request identifier
func (req *Request) Identifier() RequestIdentifier {
	return req.identifier
}

// AwaitReply blocks the calling goroutine
// until either the reply is fulfilled or failed, the request timed out
// a user-defined deadline was exceeded or the request was prematurely canceled.
// The timer is started when AwaitReply is called.
func (req *Request) AwaitReply(ctx context.Context) (webwire.Payload, error) {
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
		return &webwire.EncodedPayload{}, webwire.NewTimeoutErr(
			fmt.Errorf("timed out"),
		)

	case reply := <-req.reply:
		timeoutTimer.Stop()
		if reply.Error != nil {
			return nil, reply.Error
		}

		// Don't return nil even if the reply is empty
		// to prevent invalid memory access attempts
		// caused by forgetting to check for != nil
		if reply.Reply == nil {
			return &webwire.EncodedPayload{}, nil
		}

		return reply.Reply, nil
	}
}
