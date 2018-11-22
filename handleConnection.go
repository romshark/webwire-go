package webwire

import (
	"fmt"
	"time"
)

func (srv *server) writeConfMessage(sock Socket) error {
	writer, err := sock.GetWriter()
	if err != nil {
		return fmt.Errorf(
			"couldn't get writer for configuration message: %s",
			err,
		)
	}

	if _, err := writer.Write(srv.configMsg); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return fmt.Errorf(
				"couldn't close writer after failed conf message write: %s: %s",
				err,
				closeErr,
			)
		}
		return fmt.Errorf("couldn't write configuration message: %s", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("couldn't close writer: %s", err)
	}

	return nil
}

func (srv *server) handleConnection(
	connectionOptions ConnectionOptions,
	userAgent []byte,
	sock Socket,
) {
	// Send server configuration message
	if err := srv.writeConfMessage(sock); err != nil {
		srv.errorLog.Println("couldn't write config message: ", err)
		if closeErr := sock.Close(); closeErr != nil {
			srv.errorLog.Println("couldn't close socket: ", closeErr)
		}
		return
	}

	// Register connected client
	connection := newConnection(
		sock,
		userAgent,
		srv,
		connectionOptions,
	)

	srv.connectionsLock.Lock()
	srv.connections = append(srv.connections, connection)
	srv.connectionsLock.Unlock()

	// Call hook on successful connection
	srv.impl.OnClientConnected(connection)

	for {
		// Get a message buffer
		msg := srv.messagePool.Get()

		if !connection.IsActive() {
			msg.Close()
			connection.Close()
			srv.impl.OnClientDisconnected(connection, nil)
			break
		}

		// Await message
		if err := sock.Read(
			msg,
			time.Now().Add(srv.options.ReadTimeout), // Deadline
		); err != nil {
			msg.Close()

			if !err.IsCloseErr() {
				srv.warnLog.Printf("abnormal closure error: %s", err)
			}

			connection.Close()
			srv.impl.OnClientDisconnected(connection, err)
			break
		}

		// Parse & handle the message
		if err := srv.handleMessage(connection, msg); err != nil {
			srv.errorLog.Print("message handler failed: ", err)
		}
		msg.Close()
	}
}
