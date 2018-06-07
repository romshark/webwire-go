package webwire

import (
	"net/http"
)

// ServeHTTP will make the server listen for incoming HTTP requests
// eventually trying to upgrade them to WebSocket connections
func (srv *server) ServeHTTP(
	resp http.ResponseWriter,
	req *http.Request,
) {
	// Reject incoming connections during shutdown, pretend the server is temporarily unavailable
	srv.opsLock.Lock()
	if srv.shutdown {
		srv.opsLock.Unlock()
		http.Error(resp, "Server shutting down", http.StatusServiceUnavailable)
		return
	}
	srv.opsLock.Unlock()

	switch req.Method {
	case "OPTIONS":
		srv.impl.OnOptions(resp)
		return
	case "WEBWIRE":
		srv.handleMetadata(resp)
		return
	}

	if !srv.impl.BeforeUpgrade(resp, req) {
		return
	}

	// Establish connection
	conn, err := srv.connUpgrader.Upgrade(resp, req)
	if err != nil {
		srv.errorLog.Print("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	// Register connected client
	newClient := newClientAgent(conn, req.Header.Get("User-Agent"), srv)

	srv.clientsLock.Lock()
	srv.clients = append(srv.clients, newClient)
	srv.clientsLock.Unlock()

	// Call hook on successful connection
	srv.impl.OnClientConnected(newClient)

	for {
		// Await message
		message, err := conn.Read()
		if err != nil {
			if err.IsAbnormalCloseErr() {
				srv.warnLog.Printf("Abnormal closure error: %s", err)
			}

			newClient.unlink()
			srv.impl.OnClientDisconnected(newClient)
			return
		}

		// Parse message
		var msgObject Message
		msgTypeParsed, parserErr := msgObject.Parse(message)
		if !msgTypeParsed {
			// Couldn't determine message type, drop message
			continue
		} else if parserErr != nil {
			// Couldn't parse message, protocol error
			srv.warnLog.Println("Parser error:", parserErr)

			// Respond with an error but don't break the connection
			// because protocol errors are not critical errors
			srv.failMsg(newClient, &msgObject, ProtocolErr{})
			continue
		}

		// Handle message
		if err := srv.handleMessage(newClient, &msgObject); err != nil {
			srv.errorLog.Printf("CRITICAL FAILURE: %s", err)
			break
		}
	}
}
