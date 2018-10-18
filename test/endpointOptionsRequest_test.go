package test

import (
	"net/http"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	"github.com/stretchr/testify/require"
)

// TestEndpointOptions tests server endpoint using the OPTIONS method
func TestEndpointOptions(t *testing.T) {
	// Initialize webwire server
	server := setupServer(t, &serverImpl{}, wwr.ServerOptions{})

	// Initialize HTTP client
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	// Request metadata
	request, err := http.NewRequest("OPTIONS", server.Address(), nil)
	require.NoError(t, err)
	response, err := httpClient.Do(request)
	require.NoError(t, err)

	require.Equal(
		t,
		"WEBWIRE",
		response.Header.Get("Access-Control-Allow-Methods"),
	)
	require.Equal(t, "*", response.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "0", response.Header.Get("Content-Length"))
}
