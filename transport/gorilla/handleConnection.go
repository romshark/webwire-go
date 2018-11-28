package gorilla

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/qbeon/webwire-go/connopt"
)

func (srv *Transport) handleConnection(
	connectionOptions connopt.ConnectionOptions,
	userAgent []byte,
	conn *websocket.Conn,
) {
	conn.SetPongHandler(func(appData string) error {
		if err := conn.SetReadDeadline(
			time.Now().Add(srv.readTimeout),
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
			time.Now().Add(srv.readTimeout),
		); err != nil {
			return fmt.Errorf(
				"couldn't set read deadline in Ping handler: %s",
				err,
			)
		}
		return nil
	})

	srv.onNewConnection(
		connectionOptions,
		userAgent,
		NewConnectedSocket(conn),
	)
}
