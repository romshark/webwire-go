package client

import (
	"errors"
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

		// Abort if timed out
		if atomic.LoadUint32(&abortAwait) > 0 {
			clt.conn.Close()
			return
		}

		// Get a message buffer from the pool
		buf := clt.messageBufferPool.Get()

		// Abort if timed out
		if atomic.LoadUint32(&abortAwait) > 0 {
			clt.conn.Close()
			buf.Close()
			return
		}

		// Await the server configuration handshake response
		clt.conn.SetReadDeadline(deadline)
		if err := clt.conn.Read(buf); err != nil {
			result <- dialResult{err: fmt.Errorf("read err: %s", err.Error())}
			buf.Close()
			return
		}

		// Abort if timed out
		if atomic.LoadUint32(&abortAwait) > 0 {
			clt.conn.Close()
			buf.Close()
			return
		}

		msg := &message.Message{}

		// Parse the first incoming message
		parsedMessageType, parseErr := msg.Parse(buf.Data())
		if !parsedMessageType {
			result <- dialResult{err: errors.New(
				"unexpected message (unknown type)",
			)}
			buf.Close()
			return
		}
		if parseErr != nil {
			result <- dialResult{err: errors.New("message parser failed")}
			buf.Close()
			return
		}
		if msg.Type != message.MsgConf {
			result <- dialResult{err: fmt.Errorf(
				"unexpected message type: %d (expected server config message)",
				msg.Type,
			)}
			buf.Close()
			return
		}

		// Verify the protocol version
		if err := verifyProtocolVersion(
			msg.ServerConfiguration.MajorProtocolVersion,
			msg.ServerConfiguration.MinorProtocolVersion,
		); err != nil {
			result <- dialResult{err: err}
			buf.Close()
			return
		}

		// Finish successful dial
		clt.conn.SetReadDeadline(time.Time{})
		result <- dialResult{
			serverConfiguration: msg.ServerConfiguration,
		}
		buf.Close()
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
