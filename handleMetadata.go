package webwire

// handleMetadata handles endpoint metadata requests
import (
	"encoding/json"
	"net/http"
	"time"
)

func (srv *server) handleMetadata(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(resp).Encode(struct {
		ProtocolVersion string `json:"protocol-version"`
		ReadTimeout     uint32 `json:"read-timeout"`
	}{
		ProtocolVersion: protocolVersion,
		ReadTimeout:     uint32(srv.options.ReadTimeout / time.Second),
	})
}
