package client

import (
	"log"
	"sync"
	"time"

	"github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/message"
)

// heartbeat represents the heartbeat module of the client that's responsible
// for periodically pinging the server to keep the connection alive
type heartbeat struct {
	ticker    *time.Ticker
	resetChan chan struct{}
	stopChan  chan struct{}
	lock      sync.Mutex
	conn      webwire.Socket
	errLog    *log.Logger
}

// newHeartbeat creates a new heartbeat module instance
func newHeartbeat(conn webwire.Socket, errLog *log.Logger) heartbeat {
	return heartbeat{
		resetChan: make(chan struct{}, 1),
		stopChan:  make(chan struct{}, 1),
		lock:      sync.Mutex{},
		conn:      conn,
		errLog:    errLog,
	}
}

// start starts the heartbeating loop blocking the calling goroutine
func (hb *heartbeat) start(dur time.Duration) {
	hb.lock.Lock()
	if hb.ticker != nil {
		hb.lock.Unlock()
		return
	}
	hb.ticker = time.NewTicker(dur)
	hb.lock.Unlock()
MAINLOOP:
	for {
		select {
		case <-hb.ticker.C:
			// heartbeat
			if err := hb.conn.Write([]byte{message.MsgHeartbeat}); err != nil {
				hb.errLog.Printf("couldn't send heartbeat: %s", err)
			}
		case <-hb.resetChan:
			// reset
			hb.lock.Lock()
			hb.ticker.Stop()
			hb.ticker = time.NewTicker(dur)
			hb.lock.Unlock()
		case <-hb.stopChan:
			// stop
			hb.lock.Lock()
			hb.ticker.Stop()
			hb.ticker = nil
			hb.lock.Unlock()
			break MAINLOOP
		}
	}
}

// reset resets the timeout timer
func (hb *heartbeat) reset() {
	select {
	case hb.resetChan <- struct{}{}:
	default:
	}
}

// stop halts the heartbeating
func (hb *heartbeat) stop() {
	select {
	case hb.stopChan <- struct{}{}:
	default:
	}
}
