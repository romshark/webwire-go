package client

import (
	"sync/atomic"
	"time"
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

		// Start the reconnector goroutine if not already started.
		// If it's already started then just proceed to wait until either connected or timed out
		clt.backgroundReconnect()

		if timeout > 0 {
			// Await with timeout
			return clt.backReconn.await(timeout)
		}
		// Await indefinitely
		return clt.backReconn.await(0)
	}

	if atomic.LoadInt32(&clt.status) == StatConnected {
		return nil
	}
	if err := clt.connect(); err != nil {
		return err
	}
	return nil
}
