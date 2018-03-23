package client

import (
	"fmt"
	"net/url"
	"sync/atomic"

	"github.com/gorilla/websocket"
	webwire "github.com/qbeon/webwire-go"
)

// connect will try to establish a connection to the configured webwire server
// and try to automatically restore the session if there is any.
// If the session restoration fails connect won't fail, instead it will reset the current session
// and return normally.
// Before establishing the connection - connect verifies protocol compatibility and returns an
// error if the protocol implemented by the server doesn't match the required protocol version
// of this client instance.
func (clt *Client) connect() (err error) {
	clt.connectLock.Lock()
	defer clt.connectLock.Unlock()
	if atomic.LoadInt32(&clt.status) == StatConnected {
		return nil
	}

	if err := clt.verifyProtocolVersion(); err != nil {
		return err
	}

	connURL := url.URL{Scheme: "ws", Host: clt.serverAddr, Path: "/"}

	clt.connLock.Lock()
	clt.conn, _, err = websocket.DefaultDialer.Dial(connURL.String(), nil)
	if err != nil {
		return webwire.NewDisconnectedErr(fmt.Errorf("Dial failure: %s", err))
	}
	clt.connLock.Unlock()

	// Setup reader thread
	go func() {
		defer clt.close()
		for {
			_, message, err := clt.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(
					err,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
				) {
					// Error while reading message
					clt.errorLog.Print("Failed reading message:", err)
				}

				// Set status to disconnected if it wasn't disabled
				if atomic.LoadInt32(&clt.status) == StatConnected {
					atomic.StoreInt32(&clt.status, StatDisconnected)
				}

				// Call hook
				clt.hooks.OnDisconnected()

				// Try to reconnect if the client wasn't disabled and autoconnect is on.
				// reconnect in another goroutine to let this one die and free up the socket
				go func() {
					if clt.autoconnect && atomic.LoadInt32(&clt.status) != StatDisabled {
						if err := clt.tryAutoconnect(0); err != nil {
							clt.errorLog.Printf("Auto-reconnect failed after connection loss: %s", err)
							return
						}
					}
				}()
				return
			}
			// Try to handle the message
			if err = clt.handleMessage(message); err != nil {
				clt.warningLog.Print("Failed handling message:", err)
			}
		}
	}()

	atomic.StoreInt32(&clt.status, StatConnected)

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
		// Just log a warning and still return nil, even if session restoration failed,
		// because we only care about the connection establishment in this method
		clt.warningLog.Printf("Couldn't restore session on reconnection: %s", err)

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
