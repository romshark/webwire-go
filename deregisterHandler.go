package webwire

// deregisterHandler decrements the number of currently executed handlers
// and shuts down the server if scheduled and no more operations are left
func (srv *server) deregisterHandler(con *connection) {
	srv.opsLock.Lock()
	srv.currentOps--
	if srv.shutdown && srv.currentOps < 1 {
		close(srv.shutdownRdy)
	}
	srv.opsLock.Unlock()

	con.deregisterTask()

	// Release a handler slot
	if con.options.ConcurrencyLimit > 1 {
		con.handlerSlots.Release(1)
	}
}
