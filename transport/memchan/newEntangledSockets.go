package memchan

// NewEntangledSockets creates a new socket pair
func NewEntangledSockets(server *Transport) (srv, clt *Socket) {
	srv = newSocket(server, server.bufferSize)
	clt = newSocket(server, server.bufferSize)
	entangleSockets(srv, clt)
	return
}
