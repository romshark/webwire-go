package client

// DisconnectedErr is an error type indicating that the client isn't connected to the server
type DisconnectedErr struct{}

func (err DisconnectedErr) Error() string {
	return "Client is disconnected"
}
