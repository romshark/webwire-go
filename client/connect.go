package client

import (
	"context"
	"sync/atomic"
	"time"
)

// connect will try to establish a connection to the configured webwire server
// and try to automatically restore the session if there is any.
// If the session restoration fails connect won't fail,
// instead it will reset the current session and return normally.
// Before establishing the connection - connect verifies
// protocol compatibility and returns an error if
// the protocol implemented by the server doesn't match
// the required protocol version of this client instance.
func (clt *client) connect() error {
	clt.connectLock.Lock()
	if clt.Status() == StatusConnected {
		clt.connectLock.Unlock()
		return nil
	}

	// Dial and await approval
	endpointMeta, err := clt.dial()
	if err != nil {
		clt.connectLock.Unlock()
		return err
	}

	// Start heartbeat
	go clt.heartbeat.start(
		endpointMeta.ReadTimeout - endpointMeta.ReadTimeout/4,
	)

	// Setup reader thread
	go func() {
		defer func() {
			// Set status
			clt.setStatus(StatusDisconnected)
			select {
			case clt.readerClosing <- true:
			default:
			}
		}()
		for {
			// Get a message buffer from the pool
			msg := clt.messagePool.Get()

			if err := clt.conn.Read(msg, time.Time{}); err != nil {
				// Return message object back to the pool
				msg.Close()

				// Set connection status to disconnected
				clt.setStatus(StatusDisconnected)

				// Stop heartbeat
				clt.heartbeat.stop()

				// Call hook
				clt.impl.OnDisconnected()

				// Try to reconnect if autoconn wasn't disabled.
				// reconnect in another goroutine to let this one die
				// and free up the socket
				if atomic.LoadInt32(&clt.autoconnect) == autoconnectEnabled {
					go func() {
						clt.tryAutoconnect(context.Background(), false)
					}()
				}
				return
			}

			// Try to handle the message
			if err := clt.handleMessage(msg); err != nil {
				clt.options.ErrorLog.Print("message handler failed:", err)
			}
		}
	}()

	clt.setStatus(StatusConnected)

	// Read the current sessions key if there is any
	clt.sessionLock.RLock()
	if clt.session == nil {
		clt.sessionLock.RUnlock()
		clt.connectLock.Unlock()
		return nil
	}
	sessionKey := clt.session.Key
	clt.sessionLock.RUnlock()

	// Set session restoration deadline
	ctx, closeCtx := context.WithTimeout(
		context.Background(),
		clt.options.DefaultRequestTimeout,
	)

	// Try to restore session if necessary
	restoredSession, err := clt.requestSessionRestoration(
		ctx,
		[]byte(sessionKey),
	)
	if err != nil {
		// Just log a warning and still return nil,
		// even if session restoration failed,
		// because we only care about the connection establishment
		// in this method
		clt.options.WarnLog.Printf(
			"Couldn't restore session on reconnection: %s",
			err,
		)

		// Reset the session
		clt.sessionLock.Lock()
		clt.session = nil
		clt.sessionLock.Unlock()
		clt.connectLock.Unlock()
		closeCtx()
		return nil
	}

	clt.sessionLock.Lock()
	clt.session = restoredSession
	clt.sessionLock.Unlock()
	clt.connectLock.Unlock()
	closeCtx()
	return nil
}
