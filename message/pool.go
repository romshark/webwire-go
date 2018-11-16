package message

// Pool defines the message buffer pool interface
type Pool interface {
	// Get fetches a message buffer from the pool which must be put back when
	// it's no longer needed
	Get() *Message
}
