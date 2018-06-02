package client

import (
	"log"
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
	log.Print("BACKGROUND RECONNECT")
	go func() {
		for {
			log.Print("TRY CONNECT")
			err := clt.connect()
			switch err := err.(type) {
			case nil:
				log.Print("RECONNECTED")
				clt.connectingLock.Lock()
				clt.backReconn.flush(nil)
				clt.connecting = false
				clt.connectingLock.Unlock()
				return
			case webwire.DisconnectedErr:
				log.Print("WAIT AND RETRY")
				time.Sleep(clt.reconnInterval)
			default:
				// Unexpected error
				clt.backReconn.flush(err)
				return
			}
		}
	}()
}
