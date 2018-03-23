package main

import (
	"context"
	"fmt"

	wwr "github.com/qbeon/webwire-go"
)

// onRequest handles incoming client requests and dispatches them to corresponding handlers
func onRequest(ctx context.Context) (wwr.Payload, error) {
	msg := ctx.Value(wwr.Msg).(wwr.Message)

	switch msg.Name {
	case "auth":
		return onAuth(ctx)
	case "msg":
		return onMessage(ctx)
	}
	return wwr.Payload{}, wwr.ReqErr{
		Code:    "BAD_REQUEST",
		Message: fmt.Sprintf("Unsupported request name: %s", msg.Name),
	}
}
