package client

import (
	"context"
	"fmt"
)

// RestoreSession tries to restore the previously opened session.
// Fails if a session is currently already active
func (clt *client) RestoreSession(
	ctx context.Context,
	sessionKey []byte,
) error {
	if ctx == nil {
		ctx = context.Background()
	}

	// Apply exclusive lock
	clt.apiLock.Lock()

	// Set default deadline if no deadline is yet specified
	closeCtx := func() {}
	_, deadlineWasSet := ctx.Deadline()
	if !deadlineWasSet {
		ctx, closeCtx = context.WithTimeout(
			ctx,
			clt.options.DefaultRequestTimeout,
		)
	}

	clt.sessionLock.RLock()
	if clt.session != nil {
		clt.apiLock.Unlock()
		clt.sessionLock.RUnlock()
		closeCtx()
		return fmt.Errorf(
			"Can't restore session if another one is already active",
		)
	}
	clt.sessionLock.RUnlock()

	if err := clt.tryAutoconnect(ctx, deadlineWasSet); err != nil {
		clt.apiLock.Unlock()
		closeCtx()
		return err
	}

	restoredSession, err := clt.requestSessionRestoration(ctx, sessionKey)
	if err != nil {
		clt.apiLock.Unlock()
		closeCtx()
		return err
	}

	clt.sessionLock.Lock()
	clt.session = restoredSession
	clt.sessionLock.Unlock()

	clt.apiLock.Unlock()
	closeCtx()
	return nil
}
