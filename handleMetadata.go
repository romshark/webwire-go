package webwire

// handleMetadata handles endpoint metadata requests
import (
	"encoding/json"
	"net/http"
)

func (srv *server) handleMetadata(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(resp).Encode(struct {
		ProtocolVersion string `json:"protocol-version"`
	}{
		protocolVersion,
	})
}
