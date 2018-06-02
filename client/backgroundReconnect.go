package client

import (
	"time"

	webwire "github.com/qbeon/webwire-go"
)

func (clt *Client) backgroundReconnect() {
	clt.connectingLock.Lock()
	defer clt.connectingLock.Unlock()
	if clt.connecting {
		return
	}
	clt.connecting = true
	go func() {
		for {
			err := clt.connect()
			switch err := err.(type) {
			case nil:
				clt.connectingLock.Lock()
				clt.backReconn.flush(nil)
				clt.connecting = false
				clt.connectingLock.Unlock()
				return
			case webwire.DisconnectedErr:
				time.Sleep(clt.reconnInterval)
			default:
				// Unexpected error
				clt.backReconn.flush(err)
				return
			}
		}
	}()
}
