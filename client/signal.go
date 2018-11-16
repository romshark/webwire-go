package client

import (
	"context"
	"fmt"

	webwire "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/wwrerr"
)

// Signal sends a signal containing the given payload to the server
func (clt *client) Signal(
	ctx context.Context,
	name []byte,
	pld webwire.Payload,
) error {
	// Apply shared lock
	clt.apiLock.RLock()

	// Set default deadline if no deadline is yet specified
	closeCtx := func() {}
	_, deadlineIsSet := ctx.Deadline()
	if !deadlineIsSet {
		ctx, closeCtx = context.WithTimeout(
			ctx,
			clt.options.DefaultRequestTimeout,
		)
	}

	if err := clt.tryAutoconnect(ctx, deadlineIsSet); err != nil {
		clt.apiLock.RUnlock()

		closeCtx()
		return err
	}

	// Require either a name or a payload or both
	if len(name) < 1 && len(pld.Data) < 1 {
		clt.apiLock.RUnlock()

		closeCtx()
		return wwrerr.ProtocolErr{
			Cause: fmt.Errorf("Invalid request, request message requires " +
				"either a name, a payload or both but is missing both",
			),
		}
	}

	writer, err := clt.conn.GetWriter()
	if err != nil {
		closeCtx()
		return err
	}

	if err := message.WriteMsgSignal(
		writer,
		name,
		pld.Encoding,
		pld.Data,
		true,
	); err != nil {
		closeCtx()
		return err
	}

	clt.heartbeat.reset()

	clt.apiLock.RUnlock()

	closeCtx()
	return nil
}
