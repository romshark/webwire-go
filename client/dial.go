package client

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/qbeon/webwire-go/message"
	"github.com/qbeon/webwire-go/wwrerr"
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
		if err := clt.conn.Dial(serverAddr, deadline); err != nil {
			result <- dialResult{err: err}
			return
		}

		// Abort if timed out
		if atomic.LoadUint32(&abortAwait) > 0 {
			clt.conn.Close()
			return
		}

		// Get a message buffer from the pool
		msg := clt.messagePool.Get()

		// Abort if timed out
		if atomic.LoadUint32(&abortAwait) > 0 {
			clt.conn.Close()
			msg.Close()
			return
		}

		// Await the server configuration handshake response
		if err := clt.conn.Read(msg, deadline); err != nil {
			result <- dialResult{err: fmt.Errorf("read err: %s", err.Error())}
			clt.conn.Close()
			msg.Close()
			return
		}

		// Abort if timed out
		if atomic.LoadUint32(&abortAwait) > 0 {
			clt.conn.Close()
			msg.Close()
			return
		}

		if msg.MsgType != message.MsgConf {
			result <- dialResult{err: fmt.Errorf(
				"unexpected message type: %d (expected server config message)",
				msg.MsgType,
			)}
			clt.conn.Close()
			msg.Close()
			return
		}

		// Verify the protocol version
		if err := verifyProtocolVersion(
			msg.ServerConfiguration.MajorProtocolVersion,
			msg.ServerConfiguration.MinorProtocolVersion,
		); err != nil {
			result <- dialResult{err: err}
			msg.Close()
			return
		}

		// Finish successful dial
		result <- dialResult{
			serverConfiguration: msg.ServerConfiguration,
		}
		msg.Close()
	}()

	select {
	case <-dialingTimer.C:
		// Abort due to timeout
		atomic.StoreUint32(&abortAwait, 1)
		return message.ServerConfiguration{}, wwrerr.DialTimeoutErr{}
	case result := <-result:
		if !dialingTimer.Stop() {
			<-dialingTimer.C
		}
		return result.serverConfiguration, result.err
	}
}
