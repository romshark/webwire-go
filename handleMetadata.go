package webwire

// handleMetadata handles endpoint metadata requests
import (
	"encoding/json"
	"time"

	"github.com/valyala/fasthttp"
)

func (srv *server) handleMetadata(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(ctx).Encode(struct {
		ProtocolVersion string `json:"protocol-version"`
		ReadTimeout     uint32 `json:"read-timeout"`
	}{
		ProtocolVersion: protocolVersion,
		ReadTimeout:     uint32(srv.options.ReadTimeout / time.Second),
	}); err != nil {
		panic(err)
	}
}
