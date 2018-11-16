package shared

// ChatMessage represents a chat message containing the senders name
type ChatMessage struct {
	User string `json:"user"`
	Msg  []byte `json:"msg"`
}
