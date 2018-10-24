package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/qbeon/webwire-go"
)

type endpointMeta struct {
	ReadTimeout time.Duration
}

// verifyProtocolVersion requests the endpoint metadata
// to verify the server is running a supported protocol version
func (clt *client) verifyProtocolVersion() (endpointMeta, error) {
	// Clone TLS configuration (if any)
	var tlsConfig *tls.Config
	if clt.tlsConfig != nil {
		tlsConfig = clt.tlsConfig.Clone()
	}

	// Initialize HTTP client
	var httpClient = &http.Client{
		Timeout: clt.reconnInterval,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: clt.reconnInterval,
			}).Dial,
			TLSHandshakeTimeout: clt.reconnInterval,
			TLSClientConfig:     tlsConfig,
		},
	}

	addr := clt.serverAddr.String()
	request, err := http.NewRequest("WEBWIRE", addr, nil)
	if err != nil {
		// Panic because the request is always expected to be valid
		// except something is broken here internally
		panic(fmt.Errorf("Couldn't create HTTP metadata request: %s", err))
	}
	response, err := httpClient.Do(request)
	if err != nil {
		return endpointMeta{}, webwire.NewDisconnectedErr(fmt.Errorf(
			"Endpoint metadata request failed: %s", err,
		))
	}

	// Read response body
	defer response.Body.Close()
	encodedData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return endpointMeta{}, webwire.NewProtocolErr(fmt.Errorf(
			"Couldn't read metadata response body: %s",
			err,
		))
	}

	if response.StatusCode == http.StatusServiceUnavailable {
		return endpointMeta{}, webwire.NewDisconnectedErr(fmt.Errorf(
			"Endpoint unavailable: %s",
			response.Status,
		))
	}

	// Unmarshal response
	var metadata struct {
		ProtocolVersion string `json:"protocol-version"`
		ReadTimeout     uint32 `json:"read-timeout"`
	}
	if err := json.Unmarshal(encodedData, &metadata); err != nil {
		return endpointMeta{}, webwire.NewProtocolErr(fmt.Errorf(
			"Couldn't parse HTTP metadata response ('%s'): %s",
			string(encodedData),
			err,
		))
	}

	// Verify metadata
	if metadata.ProtocolVersion != supportedProtocolVersion {
		return endpointMeta{}, webwire.NewConnIncompErr(
			metadata.ProtocolVersion,
			supportedProtocolVersion,
		)
	}

	return endpointMeta{
		ReadTimeout: time.Duration(metadata.ReadTimeout) * time.Second,
	}, nil
}
