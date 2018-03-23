package client

import (
	"fmt"
	"sync/atomic"
	"time"

	webwire "github.com/qbeon/webwire-go"
)

func (clt *Client) tryAutoconnect(timeout time.Duration) error {
	fmt.Println("TRY AUTOCONNECT")
	if atomic.LoadInt32(&clt.status) == StatConnected {
		fmt.Println("ALREADY CONNECTED")
		return nil
	}

	if clt.autoconnect {
		stopTrying := make(chan error, 1)
		connected := make(chan error, 1)
		go func() {
			fmt.Println("AUTOCONNECT")
			for {
				select {
				case <-stopTrying:
					return
				default:
				}

				err := clt.connect()
				switch err := err.(type) {
				case nil:
					fmt.Println("CONNECTED!")
					close(connected)
					return
				case webwire.DisconnectedErr:
					fmt.Println("RETRY")
					time.Sleep(clt.reconnInterval)
				default:
					// Unexpected error
					connected <- err
					return
				}
			}
		}()

		if timeout > 0 {
			select {
			case err := <-connected:
				return err
			case <-time.After(timeout):
				// Stop reconnection trial loop and return timeout error
				close(stopTrying)
				// TODO: set missing timeout target
				return webwire.ReqTimeoutErr{}
			}
		} else {
			// Try indefinitely
			return <-connected
		}
	} else {
		if err := clt.connect(); err != nil {
			return err
		}
	}
	return nil
}
