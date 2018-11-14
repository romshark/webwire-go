package memchan

// entangleSockets connects two sockets
func entangleSockets(server, client *Socket) {
	if server.status != nil {
		panic("the server socket is already entangled")
	}
	if client.status != nil {
		panic("the server socket is already entangled")
	}

	// Set the socket types
	server.sockType = SocketServer
	client.sockType = SocketClient

	// Entangle references
	server.remote = client
	client.remote = server

	// Initialize shared status
	status := statusDisconnected
	server.status = &status
	client.status = &status
}
