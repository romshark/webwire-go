package client

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/qbeon/webwire-go/message"
)

// dial tries to dial in the server and await an approval including the endpoint
// metadata before the configured dialing timeout is reached.
// clt.dial should only be called from within clt.connect.
func (clt *client) dial() (message.ServerConfiguration, error) {
	dialingTimer := time.NewTimer(clt.options.DialingTimeout)
	deadline := time.Now().Add(clt.options.DialingTimeout)

	type dialResult struct {
		serverConfiguration message.ServerConfiguration
		err                 error
	}

	result := make(chan dialResult, 1)
	abortAwait := uint32(0)

	serverAddr := clt.serverAddr
	if serverAddr.Scheme == "https" {
		serverAddr.Scheme = "wss"
	} else {
		serverAddr.Scheme = "ws"
	}

	go func() {
		// Dial
		if err := clt.conn.Dial(serverAddr); err != nil {
			result <- dialResult{err: err}
			return
		}
		// Close the connection if the dial succeeded after the timeout
		if atomic.LoadUint32(&abortAwait) > 0 {
			clt.conn.Close()
			return
		}
		clt.conn.SetReadDeadline(deadline)
		rawMsg, err := clt.conn.Read()
		if err != nil {
			result <- dialResult{err: fmt.Errorf("read err: %s", err.Error())}
			return
		}
		var msg message.Message
		parsedMessageType, parseErr := msg.Parse(rawMsg)
		if !parsedMessageType {
			result <- dialResult{err: fmt.Errorf(
				"unexpected message (unknown type)",
			)}
			return
		}
		if parseErr != nil {
			result <- dialResult{err: fmt.Errorf("message parser failed")}
			return
		}
		if msg.Type != message.MsgConf {
			result <- dialResult{err: fmt.Errorf(
				"unexpected message type: %d",
				msg.Type,
			)}
			return
		}
		if !verifyProtocolVersion(
			msg.ServerConfiguration.MajorProtocolVersion,
			msg.ServerConfiguration.MinorProtocolVersion,
		) {
			result <- dialResult{err: fmt.Errorf(
				"unexpected message type: %d",
				msg.Type,
			)}
			return
		}
		clt.conn.SetReadDeadline(time.Time{})
		result <- dialResult{
			serverConfiguration: msg.ServerConfiguration,
		}
	}()

	select {
	case <-dialingTimer.C:
		// Abort due to timeout
		dialingTimer.Stop()
		atomic.StoreUint32(&abortAwait, 1)
		return message.ServerConfiguration{}, DialTimeout{}
	case result := <-result:
		dialingTimer.Stop()
		return result.serverConfiguration, result.err
	}
}
