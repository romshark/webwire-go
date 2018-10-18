package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	"github.com/stretchr/testify/require"
)

// TestEndpointMetadata tests server endpoint metadata
func TestEndpointMetadata(t *testing.T) {
	expectedVersion := "1.4"

	// Initialize webwire server
	server := setupServer(t, &serverImpl{}, wwr.ServerOptions{})

	// Initialize HTTP client
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	// Request metadata
	request, err := http.NewRequest("WEBWIRE", server.Address(), nil)
	require.NoError(t, err)
	response, err := httpClient.Do(request)
	require.NoError(t, err)

	// Read response body
	defer response.Body.Close()
	encodedData, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	// Unmarshal response
	var metadata struct {
		ProtocolVersion string `json:"protocol-version"`
	}
	require.NoError(t, json.Unmarshal(encodedData, &metadata))

	// Verify metadata
	require.Equal(t, expectedVersion, metadata.ProtocolVersion)
}
