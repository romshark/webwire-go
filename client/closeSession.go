package client

import (
	"context"
	"sync/atomic"

	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/payload"
)

// CloseSession disables the currently active session
// and acknowledges the server if connected.
// The session will be destroyed if this is it's last connection remaining.
// If the client is not connected then the synchronization is skipped.
// Does nothing if there's no active session
func (clt *client) CloseSession() error {
	// Apply exclusive lock
	clt.apiLock.Lock()

	clt.sessionLock.RLock()
	if clt.session == nil {
		clt.sessionLock.RUnlock()
		clt.apiLock.Unlock()
		return nil
	}
	clt.sessionLock.RUnlock()

	// Synchronize session closure to the server if connected
	if atomic.LoadInt32(&clt.status) == Connected {
		if _, err := clt.sendNamelessRequest(
			context.Background(),
			message.MsgCloseSession,
			payload.Payload{},
			clt.options.DefaultRequestTimeout,
		); err != nil {
			clt.apiLock.Unlock()
			return err
		}
	}

	// Reset session locally after destroying it on the server
	clt.sessionLock.Lock()
	clt.session = nil
	clt.sessionLock.Unlock()

	clt.apiLock.Unlock()
	return nil
}
