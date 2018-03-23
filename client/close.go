package client

import (
	"sync/atomic"
)

func (clt *Client) close() {
	clt.connLock.Lock()
	defer clt.connLock.Unlock()
	if atomic.LoadInt32(&clt.status) < StatConnected {
		// Either disconnected or disabled
		return
	}
	if err := clt.conn.Close(); err != nil {
		clt.errorLog.Printf("Failed closing connection: %s", err)
	}
	atomic.StoreInt32(&clt.status, StatDisabled)
}
