package client

import (
	"context"

	webwire "github.com/qbeon/webwire-go"
)

// Request sends a request containing the given payload to the server
// and asynchronously returns the servers response
// blocking the calling goroutine.
// Returns an error if the request failed for some reason
func (clt *client) Request(
	ctx context.Context,
	name []byte,
	payload webwire.Payload,
) (webwire.Reply, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Apply shared lock
	clt.apiLock.RLock()

	// Set default deadline if no deadline is yet specified
	closeCtx := func() {}
	_, deadlineWasSet := ctx.Deadline()
	if !deadlineWasSet {
		ctx, closeCtx = context.WithTimeout(
			ctx,
			clt.options.DefaultRequestTimeout,
		)
	}

	if err := clt.tryAutoconnect(ctx, deadlineWasSet); err != nil {
		clt.apiLock.RUnlock()

		closeCtx()
		return nil, err
	}

	reply, err := clt.sendRequest(
		ctx,
		determineMsgTypeBasedOnEncoding(payload.Encoding),
		name,
		payload,
		clt.options.DefaultRequestTimeout,
	)

	clt.apiLock.RUnlock()
	closeCtx()

	return reply, err
}
