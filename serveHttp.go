package webwire

import (
	"fmt"
	"net/http"
	"time"
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

	// Set ping/pong handlers
	conn.OnPong(func(string) error {
		if err := conn.SetReadDeadline(
			time.Now().Add(srv.options.HearthbeatTimeout),
		); err != nil {
			return fmt.Errorf(
				"Couldn't set read deadline in Pong handler: %s",
				err,
			)
		}
		return nil
	})
	conn.OnPing(func(string) error {
		if err := conn.SetReadDeadline(
			time.Now().Add(srv.options.HearthbeatTimeout),
		); err != nil {
			return fmt.Errorf(
				"Couldn't set read deadline in Ping handler: %s",
				err,
			)
		}
		return nil
	})
	if err := conn.SetReadDeadline(
		time.Now().Add(srv.options.HearthbeatTimeout),
	); err != nil {
		srv.errorLog.Printf("Couldn't set read deadline: %s", err)
		return
	}

	// Register connected client
	newClient := newClientAgent(conn, req.Header.Get("User-Agent"), srv)

	srv.clientsLock.Lock()
	srv.clients = append(srv.clients, newClient)
	srv.clientsLock.Unlock()

	// Call hook on successful connection
	srv.impl.OnClientConnected(newClient)

	// Start hearthbeat sender
	stopHearthbeating := make(chan struct{}, 1)
	go func() {
		hearthbeatTicker := time.NewTicker(srv.options.HearthbeatInterval)
	HEARTHBEAT_LOOP:
		for {
			if err := conn.WritePing(
				nil,
				time.Now().Add(srv.options.HearthbeatInterval),
			); err != nil {
				srv.errorLog.Printf("Couldn't write ping frame: %s", err)
			}
			select {
			case <-hearthbeatTicker.C:
				// Just continue
			case <-stopHearthbeating:
				hearthbeatTicker.Stop()
				break HEARTHBEAT_LOOP
			}
		}
	}()

	for {
		// Await message
		message, err := conn.Read()
		if err != nil {
			if err.IsAbnormalCloseErr() {
				srv.warnLog.Printf("Abnormal closure error: %s", err)
			}

			newClient.unlink()
			srv.impl.OnClientDisconnected(newClient)
			break
		}

		// Parse & handle the message
		go srv.handleMessage(newClient, message)
	}

	// Connection closed
	stopHearthbeating <- struct{}{}
}
