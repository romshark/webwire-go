package client

import (
	"sync/atomic"
	"time"

	webwire "github.com/qbeon/webwire-go"
)

func (clt *Client) tryAutoconnect(timeout time.Duration) error {
	if atomic.LoadInt32(&clt.isConnected) > 0 {
		return nil
	}

	if clt.autoconnect {
		stopTrying := make(chan error, 1)
		connected := make(chan error, 1)
		go func() {
			for {
				select {
				case <-stopTrying:
					return
				default:
				}

				err := clt.connect()
				switch err := err.(type) {
				case nil:
					close(connected)
					return
				case webwire.DisconnectedErr:
					time.Sleep(clt.reconnInterval)
				default:
					// Unexpected error
					connected <- err
					return
				}
			}
		}()

		// TODO: implement autoconnect
		select {
		case err := <-connected:
			return err
		case <-time.After(timeout):
			// Stop reconnection trial loop and return timeout error
			close(stopTrying)
			return webwire.ReqTimeoutErr{}
		}
	} else {
		if err := clt.connect(); err != nil {
			return err
		}
	}
	return nil
}
