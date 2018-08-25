package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/qbeon/webwire-go"
)

// verifyProtocolVersion requests the endpoint metadata
// to verify the server is running a supported protocol version
func (clt *client) verifyProtocolVersion() error {
	// Initialize HTTP client
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	request, err := http.NewRequest(
		"WEBWIRE", "http://"+clt.serverAddr+"/", nil,
	)
	if err != nil {
		panic(fmt.Errorf("Couldn't create HTTP metadata request: %s", err))
	}
	response, err := httpClient.Do(request)
	if err != nil {
		return webwire.NewDisconnectedErr(fmt.Errorf(
			"Endpoint metadata request failed: %s", err,
		))
	}

	// Read response body
	defer response.Body.Close()
	encodedData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return webwire.NewProtocolErr(fmt.Errorf(
			"Couldn't read metadata response body: %s",
			err,
		))
	}

	if response.StatusCode == http.StatusServiceUnavailable {
		return webwire.NewDisconnectedErr(fmt.Errorf(
			"Endpoint unavailable: %s",
			response.Status,
		))
	}

	// Unmarshal response
	var metadata struct {
		ProtocolVersion string `json:"protocol-version"`
	}
	if err := json.Unmarshal(encodedData, &metadata); err != nil {
		return webwire.NewProtocolErr(fmt.Errorf(
			"Couldn't parse HTTP metadata response ('%s'): %s",
			string(encodedData),
			err,
		))
	}

	// Verify metadata
	if metadata.ProtocolVersion != supportedProtocolVersion {
		return webwire.NewConnIncompErr(
			metadata.ProtocolVersion,
			supportedProtocolVersion,
		)
	}

	return nil
}
