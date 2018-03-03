package client

import "sync/atomic"

func (clt *Client) close() {
	clt.connLock.Lock()
	defer clt.connLock.Unlock()
	if atomic.LoadInt32(&clt.isConnected) < 1 {
		return
	}
	if err := clt.conn.Close(); err != nil {
		clt.errorLog.Printf("Failed closing connection: %s", err)
	}
	atomic.StoreInt32(&clt.isConnected, 0)
}
