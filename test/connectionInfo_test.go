package test

import (
	"regexp"
	"testing"
	"time"

	tmdwg "github.com/qbeon/tmdwg-go"
	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConnectionInfo tests the connection.Info method
func TestConnectionInfo(t *testing.T) {
	handlerFinished := tmdwg.NewTimedWaitGroup(1, 1*time.Second)

	// Initialize server
	setup := setupTestServer(
		t,
		&serverImpl{
			onClientConnected: func(conn wwr.Connection) {
				info := conn.Info()
				assert.WithinDuration(
					t,
					time.Now(),
					info.ConnectionTime,
					1*time.Second,
				)

				switch *argTransport {
				case "fasthttp/websocket":
					// Check user agent string
					match, err := regexp.Match(
						"^Go-http-client/1\\.1$",
						info.UserAgent,
					)
					assert.NoError(t, err)
					assert.True(t, match)

					// Check remote address
					assert.NotNil(t, info.RemoteAddr)

				case "memchan":
					// Check user agent string
					match, err := regexp.Match(
						"^webwire memchan client \\(0x[A-Fa-f0-9]{6,12}\\)$",
						info.UserAgent,
					)
					assert.NoError(t, err)
					assert.True(t, match)

					// Check remote address
					assert.NotNil(t, info.RemoteAddr)

				default:
					t.Fatalf(
						"unexpected server transport implementation: %s",
						*argTransport,
					)
				}

				handlerFinished.Progress(1)
			},
		},
		wwr.ServerOptions{},
		nil, // Use the default transport implementation
	)

	// Initialize client
	setup.newClient(
		wwrclt.Options{},
		nil, // Use the default transport implementation
		testClientHooks{},
	)

	require.NoError(t, handlerFinished.Wait())
}
