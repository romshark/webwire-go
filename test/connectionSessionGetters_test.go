package test

import (
	"context"
	"sync"
	"testing"
	"time"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/payload"
	"github.com/stretchr/testify/assert"
)

type testSessInfo struct {
	UserIdent  string
	SomeNumber int
}

// Copy implements the webwire.SessionInfo interface
func (sinf *testSessInfo) Copy() wwr.SessionInfo {
	return &testSessInfo{
		UserIdent:  sinf.UserIdent,
		SomeNumber: sinf.SomeNumber,
	}
}

// Fields implements the webwire.SessionInfo interface
func (sinf *testSessInfo) Fields() []string {
	return []string{"uid", "some-number"}
}

// Copy implements the webwire.SessionInfo interface
func (sinf *testSessInfo) Value(fieldName string) interface{} {
	switch fieldName {
	case "uid":
		return sinf.UserIdent
	case "some-number":
		return sinf.SomeNumber
	}
	return nil
}

// TestConnectionSessionGetters tests the connection session information getters
func TestConnectionSessionGetters(t *testing.T) {
	signalDone := sync.WaitGroup{}
	signalDone.Add(1)

	compareSession := func(conn wwr.Connection) {
		timeNow := time.Now()

		sess := conn.Session()
		assert.NotNil(t, sess)

		assert.Equal(t, "testsessionkey", sess.Key)
		assert.Equal(t, &testSessInfo{
			UserIdent:  "clientidentifiergoeshere", // uid
			SomeNumber: 12345,                      // some-number
		}, sess.Info)
		assert.WithinDuration(t, timeNow, sess.Creation, 1*time.Second)
		assert.WithinDuration(t, timeNow, sess.LastLookup, 1*time.Second)

		assert.WithinDuration(
			t,
			timeNow,
			conn.SessionCreation(),
			1*time.Second,
		)
		assert.Equal(t, "testsessionkey", conn.SessionKey())
		uid := conn.SessionInfo("uid")
		assert.NotNil(t, uid)
		assert.IsType(t, string(""), uid)

		someNumber := conn.SessionInfo("some-number")
		assert.NotNil(t, someNumber)
		assert.IsType(t, int(0), someNumber)
	}

	// Initialize server
	setup := SetupTestServer(
		t,
		&ServerImpl{
			ClientConnected: func(_ wwr.ConnectionOptions, c wwr.Connection) {
				// Before session creation
				assert.Nil(t, c.Session())
				assert.True(t, c.SessionCreation().IsZero())
				assert.Equal(t, "", c.SessionKey())
				assert.Nil(t, c.SessionInfo("uid"))
				assert.Nil(t, c.SessionInfo("some-number"))

				assert.NoError(t, c.CreateSession(
					&testSessInfo{
						UserIdent:  "clientidentifiergoeshere", // uid
						SomeNumber: 12345,                      // some-number
					},
				))

				// After session creation
				compareSession(c)
			},
			Request: func(
				_ context.Context,
				c wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				compareSession(c)
				return wwr.Payload{}, nil
			},
			Signal: func(_ context.Context, c wwr.Connection, _ wwr.Message) {
				defer signalDone.Done()
				compareSession(c)
			},
		},
		wwr.ServerOptions{
			SessionKeyGenerator: &SessionKeyGen{
				OnGenerate: func() string {
					return "testsessionkey"
				},
			},
		},
		nil, // Use the default transport implementation
	)

	// Connect new client
	sock, _ := setup.NewClientSocket()

	readSessionCreated(t, sock)

	requestSuccess(t, sock, 32, []byte("r"), payload.Payload{})

	signal(t, sock, []byte("s"), payload.Payload{})

	signalDone.Wait()
}
