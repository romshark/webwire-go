package gorilla

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// SockReadErr implements the SockReadErr interface
type SockReadErr struct {
	cause error
}

// Error implements the Go error interface
func (err SockReadErr) Error() string {
	return fmt.Sprintf("reading socket failed: %s", err.cause)
}

// IsCloseErr implements the SockReadErr interface
func (err SockReadErr) IsCloseErr() bool {
	return websocket.IsCloseError(
		err.cause,
		websocket.CloseNormalClosure,
		websocket.CloseGoingAway,
		websocket.CloseAbnormalClosure,
	)
}

// SockReadWrongMsgTypeErr implements the SockReadErr interface
type SockReadWrongMsgTypeErr struct {
	messageType int
}

// Error implements the Go error interface
func (err SockReadWrongMsgTypeErr) Error() string {
	return fmt.Sprintf("invalid websocket message type: %d", err.messageType)
}

// IsCloseErr implements the SockReadErr interface
func (err SockReadWrongMsgTypeErr) IsCloseErr() bool {
	return false
}
