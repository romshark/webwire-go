package client

import (
	"fmt"
	"net/url"
	"sync/atomic"

	"github.com/gorilla/websocket"
	webwire "github.com/qbeon/webwire-go"
)

func (clt *Client) connect() (err error) {
	clt.connectLock.Lock()
	defer clt.connectLock.Unlock()
	if atomic.LoadInt32(&clt.isConnected) > 0 {
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
					break
				} else {
					// Shutdown client due to clean disconnection
					break
				}
			}
			// Try to handle the message
			if err = clt.handleMessage(message); err != nil {
				clt.warningLog.Print("Failed handling message:", err)
			}
		}
	}()

	atomic.StoreInt32(&clt.isConnected, 1)

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
