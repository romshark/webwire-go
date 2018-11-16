package webwire

// isShuttingDown returns true if the server is currently shutting down,
// otherwise returns false
func (srv *server) isShuttingDown() bool {
	srv.opsLock.Lock()
	if srv.shutdown {
		srv.opsLock.Unlock()
		return true
	}
	srv.opsLock.Unlock()
	return false
}
