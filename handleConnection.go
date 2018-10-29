package webwire

import (
	"fmt"
	"time"

	"github.com/fasthttp/websocket"
)

func (srv *server) handleConnection(conn *websocket.Conn) {
	conn.SetPongHandler(func(appData string) error {
		if err := conn.SetReadDeadline(
			time.Now().Add(srv.options.ReadTimeout),
		); err != nil {
			return fmt.Errorf(
				"couldn't set read deadline in Pong handler: %s",
				err,
			)
		}
		return nil
	})

	conn.SetPingHandler(func(appData string) error {
		if err := conn.SetReadDeadline(
			time.Now().Add(srv.options.ReadTimeout),
		); err != nil {
			return fmt.Errorf(
				"couldn't set read deadline in Ping handler: %s",
				err,
			)
		}
		return nil
	})

	// Send server configuration message
	if err := conn.WriteMessage(
		websocket.BinaryMessage,
		srv.configMsg,
	); err != nil {
		if err := conn.Close(); err != nil {
			srv.errorLog.Printf(
				"couldn't close connection after failed conf msg transmission: %s",
				err,
			)
		}
		return
	}

	sock := newFasthttpConnectedSocket(conn)

	if err := conn.SetReadDeadline(
		time.Now().Add(srv.options.ReadTimeout),
	); err != nil {
		srv.errorLog.Printf("couldn't set read deadline: %s", err)
		return
	}

	// TODO: use correct user agent string and connection options
	// Register connected client
	connection := newConnection(
		sock,
		//req.Header.Get("User-Agent"),
		"Test User Agent",
		srv,
		//connectionOptions,
		ConnectionOptions{
			ConcurrencyLimit: 0,
		},
	)

	srv.connectionsLock.Lock()
	srv.connections = append(srv.connections, connection)
	srv.connectionsLock.Unlock()

	// Call hook on successful connection
	srv.impl.OnClientConnected(connection)

	for {
		// Await message
		message, err := sock.Read()
		if err != nil {
			if err.IsAbnormalCloseErr() {
				srv.warnLog.Printf("abnormal closure error: %s", err)
			}

			connection.Close()
			srv.impl.OnClientDisconnected(connection)
			break
		}

		// Parse & handle the message
		srv.handleMessage(connection, message)
	}
}
