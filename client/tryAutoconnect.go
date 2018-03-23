package client

import (
	"sync/atomic"
	"time"

	webwire "github.com/qbeon/webwire-go"
)

func (clt *Client) tryAutoconnect(timeout time.Duration) error {
	// If autoconnect is enabled the client will spawn a new autoconnector goroutine which
	// will periodically poll the server and check whether it's available again.
	// If the autoconnector goroutine has already been spawned then tryAutoconnect will
	// just await the connection or timeout respectively
	if clt.autoconnect {
		if atomic.LoadInt32(&clt.status) == StatConnected {
			return nil
		}

		clt.connectingLock.RLock()
		// Start the reconnector goroutine if not already started.
		// If it's already started then just proceed to wait until either connected or timed out
		if clt.connecting == nil {
			clt.connecting = make(chan error, 1)
			go func() {
				for {
					err := clt.connect()
					switch err := err.(type) {
					case nil:
						clt.connectingLock.Lock()
						close(clt.connecting)
						clt.connecting = nil
						clt.connectingLock.Unlock()
						return
					case webwire.DisconnectedErr:
						time.Sleep(clt.reconnInterval)
					default:
						// Unexpected error
						clt.connecting <- err
						return
					}
				}
			}()
		}
		clt.connectingLock.RUnlock()

		if timeout > 0 {
			select {
			case err := <-clt.connecting:
				return err
			case <-time.After(timeout):
				// TODO: set missing timeout target
				return webwire.ReqTimeoutErr{}
			}
		} else {
			// Try indefinitely
			return <-clt.connecting
		}
	} else {
		if atomic.LoadInt32(&clt.status) == StatConnected {
			return nil
		}
		if err := clt.connect(); err != nil {
			return err
		}
	}
	return nil
}
