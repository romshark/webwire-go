package webwire

import "time"

// heartbeat starts a heartbeat for the given connection
// blocking the calling goroutine until the stop channel is triggered
func (srv *server) heartbeat(conn Socket, stop chan struct{}) {
	hearthbeatTicker := time.NewTicker(srv.options.HeartbeatInterval)
HEARTBEAT_LOOP:
	for {
		if err := conn.WritePing(
			nil,
			time.Now().Add(srv.options.HeartbeatInterval),
		); err != nil {
			srv.errorLog.Printf("Couldn't write ping frame: %s", err)
		}
		select {
		case <-hearthbeatTicker.C:
			// Just continue
		case <-stop:
			hearthbeatTicker.Stop()
			break HEARTBEAT_LOOP
		}
	}
}
