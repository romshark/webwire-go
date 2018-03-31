package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
)

// TestEndpointMetadata tests server endpoint metadata
func TestEndpointMetadata(t *testing.T) {
	expectedVersion := "1.2"

	// Initialize webwire server
	server := setupServer(t, webwire.ServerOptions{})

	// Initialize HTTP client
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	// Request metadata
	request, err := http.NewRequest(
		"WEBWIRE",
		"http://"+server.Addr().String()+"/",
		nil,
	)
	if err != nil {
		t.Fatalf("Couldn't create HTTP request: %s", err)
	}
	response, err := httpClient.Do(request)
	if err != nil {
		t.Fatalf("HTTP request failed: %s", err)
	}

	// Read response body
	defer response.Body.Close()
	encodedData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Couldn't read response body: %s", err)
	}

	// Unmarshal response
	var metadata struct {
		ProtocolVersion string `json:"protocol-version"`
	}
	if err := json.Unmarshal(encodedData, &metadata); err != nil {
		t.Fatalf(
			"Couldn't parse HTTP response ('%s'): %s",
			string(encodedData),
			err,
		)
	}

	// Verify metadata
	if metadata.ProtocolVersion != expectedVersion {
		t.Fatalf("Unexpected protocol version: %s", metadata.ProtocolVersion)
	}
}
