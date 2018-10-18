package client

import (
	"context"
	"sync/atomic"
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
	defer clt.connectLock.Unlock()
	if atomic.LoadInt32(&clt.status) == Connected {
		return nil
	}

	if err := clt.verifyProtocolVersion(); err != nil {
		return err
	}

	serverAddr := clt.serverAddr
	if serverAddr.Scheme == "https" {
		serverAddr.Scheme = "wss"
	} else {
		serverAddr.Scheme = "ws"
	}

	if err := clt.conn.Dial(serverAddr); err != nil {
		return err
	}

	// Setup reader thread
	go func() {
		defer func() {
			// Set status
			atomic.StoreInt32(&clt.status, Disconnected)
			select {
			case clt.readerClosing <- true:
			default:
			}
		}()
		for {
			message, err := clt.conn.Read()
			if err != nil {
				if err.IsAbnormalCloseErr() {
					// Error while reading message
					clt.errorLog.Print("Abnormal closure error:", err)
				}

				atomic.StoreInt32(&clt.status, Disconnected)

				// Call hook
				clt.impl.OnDisconnected()

				// Try to reconnect if autoconn wasn't disabled.
				// reconnect in another goroutine to let this one die
				// and free up the socket
				if atomic.LoadInt32(&clt.autoconnect) == autoconnectEnabled {
					go func() {
						if err := clt.tryAutoconnect(
							context.Background(),
							0,
						); err != nil {
							clt.errorLog.Printf(
								"Auto-reconnect failed "+
									"after connection loss: %s",
								err,
							)
							return
						}
					}()
				}
				return
			}
			// Try to handle the message
			if err := clt.handleMessage(message); err != nil {
				clt.warningLog.Print("Failed handling message:", err)
			}
		}
	}()

	atomic.StoreInt32(&clt.status, Connected)

	// Read the current sessions key if there is any
	clt.sessionLock.RLock()
	if clt.session == nil {
		clt.sessionLock.RUnlock()
		return nil
	}
	sessionKey := clt.session.Key
	clt.sessionLock.RUnlock()

	// Try to restore session if necessary
	restoredSession, err := clt.requestSessionRestoration([]byte(sessionKey))
	if err != nil {
		// Just log a warning and still return nil,
		// even if session restoration failed,
		// because we only care about the connection establishment
		// in this method
		clt.warningLog.Printf(
			"Couldn't restore session on reconnection: %s",
			err,
		)

		// Reset the session
		clt.sessionLock.Lock()
		clt.session = nil
		clt.sessionLock.Unlock()
		return nil
	}

	clt.sessionLock.Lock()
	clt.session = restoredSession
	clt.sessionLock.Unlock()
	return nil
}
