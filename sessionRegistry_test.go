package webwire

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

// TestSessRegRegistration tests registration
func TestSessRegRegistration(t *testing.T) {
	reg := newSessionRegistry(0, nil)

	// Register connection with session
	clt := newConnection(nil, nil, nil, ConnectionOptions{})
	sess := NewSession(nil, func() string { return "testkey_A" })
	clt.session = &sess

	require.NoError(t, reg.register(clt))

	// Expect 1 active session
	require.Equal(t, 1, reg.activeSessionsNum())

	// Expect 1 active connection on session A
	require.Equal(t, 1, reg.sessionConnectionsNum("testkey_A"))
}

// TestSessRegActiveSessionsNum tests the ActiveSessionsNum method
func TestSessRegActiveSessionsNum(t *testing.T) {
	expectedSessionsNum := 2
	reg := newSessionRegistry(0, nil)

	// Register 2 connections on two separate sessions
	cltA1 := newConnection(nil, nil, nil, ConnectionOptions{})
	sessA := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA

	cltB1 := newConnection(nil, nil, nil, ConnectionOptions{})
	sessB := NewSession(nil, func() string { return "testkey_B" })
	cltB1.session = &sessB

	require.NoError(t, reg.register(cltA1))
	require.NoError(t, reg.register(cltB1))

	// Expect 2 active sessions
	require.Equal(t, expectedSessionsNum, reg.activeSessionsNum())

	// Expect 1 active connection on each session
	require.Equal(t, 1, reg.sessionConnectionsNum("testkey_A"))
	require.Equal(t, 1, reg.sessionConnectionsNum("testkey_B"))
}

// TestSessRegsessionConnectionsNum tests the sessionConnectionsNum method
func TestSessRegsessionConnectionsNum(t *testing.T) {
	expectedSessionsNum := 1
	reg := newSessionRegistry(0, nil)

	// Register first connection on session A
	cltA1 := newConnection(nil, nil, nil, ConnectionOptions{})
	sessA1 := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA1

	require.NoError(t, reg.register(cltA1))

	// Register second connection on same session
	cltA2 := newConnection(nil, nil, nil, ConnectionOptions{})
	sessA2 := NewSession(nil, func() string { return "testkey_A" })
	cltA2.session = &sessA2

	require.NoError(t, reg.register(cltA2))

	// Expect 1 active session
	require.Equal(t, expectedSessionsNum, reg.activeSessionsNum())

	// Expect 2 active connections on session A
	require.Equal(t, 2, reg.sessionConnectionsNum("testkey_A"))
}

// TestSessRegSessionMaxConns tests the Register method
// when the maximum number of concurrent connections of a session was reached
func TestSessRegSessionMaxConns(t *testing.T) {
	// Set the maximum number of concurrent session connection to 1
	reg := newSessionRegistry(1, nil)

	// Register first connection on session A
	cltA1 := newConnection(nil, nil, nil, ConnectionOptions{})
	sessA1 := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA1

	require.NoError(t, reg.register(cltA1))

	// Register first connection on session A
	cltA2 := newConnection(nil, nil, nil, ConnectionOptions{})
	sessA2 := NewSession(nil, func() string { return "testkey_A" })
	cltA2.session = &sessA2

	require.Error(t,
		reg.register(cltA1),
		"Expected register to return an error "+
			"due to the limit of concurrent connection being reached",
	)
}

// TestSessRegDeregistration tests deregistration
func TestSessRegDeregistration(t *testing.T) {
	reg := newSessionRegistry(0, nil)

	// Register 2 connections on two separate sessions
	cltA1 := newConnection(nil, nil, nil, ConnectionOptions{})
	sessA := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA

	cltB1 := newConnection(nil, nil, nil, ConnectionOptions{})
	sessB := NewSession(nil, func() string { return "testkey_B" })
	cltB1.session = &sessB

	require.NoError(t, reg.register(cltA1))
	require.NoError(t, reg.register(cltB1))

	// Expect 2 active sessions
	require.Equal(t, 2, reg.activeSessionsNum())

	// Expect 1 active connection on each session
	require.Equal(t, 1, reg.sessionConnectionsNum("testkey_A"))
	require.Equal(t, 1, reg.sessionConnectionsNum("testkey_B"))

	// Deregister first connection, expect 0 because session was removed
	require.Equal(t, 0, reg.deregister(cltA1, false))

	// Expect 1 active session after deregistration of the first connection
	require.Equal(t, 1, reg.activeSessionsNum())

	// Expect 0 active connections on deregistered session
	require.Equal(t, -1, reg.sessionConnectionsNum("testkey_A"))
	require.Equal(t, 1, reg.sessionConnectionsNum("testkey_B"))

	// Deregister second connection, expect 0 because session was removed
	require.Equal(t, 0, reg.deregister(cltB1, false))

	// Expect no active sessions after deregistration of the second connection
	require.Equal(t, 0, reg.activeSessionsNum())

	// Expect no active connections on both sessions
	require.Equal(t, -1, reg.sessionConnectionsNum("testkey_A"))
	require.Equal(t, -1, reg.sessionConnectionsNum("testkey_B"))
}

// TestSessRegDeregistrationMultiple tests deregistration of multiple
// connections of a single session
func TestSessRegDeregistrationMultiple(t *testing.T) {
	reg := newSessionRegistry(0, nil)

	// Register 2 connections on the same session
	cltA1 := newConnection(nil, nil, nil, ConnectionOptions{})
	sessA1 := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA1

	cltA2 := newConnection(nil, nil, nil, ConnectionOptions{})
	cltA2.session = &sessA1

	require.NoError(t, reg.register(cltA1))
	require.NoError(t, reg.register(cltA2))

	// Expect 1 active session
	require.Equal(t, 1, reg.activeSessionsNum())

	// Expect 2 active connections on session A
	require.Equal(t, 2, reg.sessionConnectionsNum("testkey_A"))

	// Deregister first connection, expect 1
	require.Equal(t, 1, reg.deregister(cltA1, false))

	// Still expect 1 active session
	// after deregistration of the first connection
	require.Equal(t, 1, reg.activeSessionsNum())

	// Expect 1 remaining active connection on deregistered session
	require.Equal(t, 1, reg.sessionConnectionsNum("testkey_A"))

	// Deregister second connection, expect 0 because session was removed
	require.Equal(t, 0, reg.deregister(cltA2, false))

	// Expect no active sessions after deregistration of the second connection
	require.Equal(t, 0, reg.activeSessionsNum())

	// Expect no active connections session A
	require.Equal(t, -1, reg.sessionConnectionsNum("testkey_A"))
}

// TestSessRegConcurrentAccess tests concurrent (de)registration
func TestSessRegConcurrentAccess(t *testing.T) {
	reg := newSessionRegistry(0, nil)
	connsToRegister := uint(16)
	registeredConns := make([]*connection, connsToRegister)
	var awaitRegistration sync.WaitGroup
	var awaitDeregistration sync.WaitGroup

	// Populate registered conns map
	for i := uint(0); i < connsToRegister; i++ {
		// Create a connection on session A
		clt := newConnection(nil, nil, nil, ConnectionOptions{})
		sess := NewSession(nil, func() string { return "testkey_A" })
		clt.session = &sess
		registeredConns[i] = clt

		awaitRegistration.Add(1)
		awaitDeregistration.Add(1)
	}

	// Concurrently register multiple connections from different goroutines
	for i := uint(0); i < connsToRegister; i++ {
		index := i
		go func() {
			assert.NoError(t, reg.register(registeredConns[index]))
			awaitRegistration.Done()
		}()
	}

	// Wait for all goroutines to finish before evaluating the results
	awaitRegistration.Wait()

	// Expect 1 active session
	require.Equal(t, 1, reg.activeSessionsNum())

	// Expect N active connections on session A
	actualSessCons := reg.sessionConnectionsNum("testkey_A")
	require.NotEqual(t, -1, actualSessCons, "Missing session testkey_A")
	require.Equal(t, connsToRegister, uint(actualSessCons))

	// Concurrently deregister multiple connections from different goroutines
	negativeReturnCounter := uint32(0)
	for i := uint(0); i < connsToRegister; i++ {
		index := i
		go func() {
			// Deregister one of the connections
			result := reg.deregister(registeredConns[index], false)
			if result == -1 {
				atomic.AddUint32(&negativeReturnCounter, 1)
			}
			awaitDeregistration.Done()
		}()
	}

	// Wait for all goroutines to finish before evaluating the results
	awaitDeregistration.Wait()

	require.Equal(t, uint32(0), negativeReturnCounter)

	// Expect 1 active session
	require.Equal(t, 0, reg.activeSessionsNum())

	// Expect no active connections on session A
	require.Equal(t, -1, reg.sessionConnectionsNum("testkey_A"))
}

// TestSessRegSessionConnections tests the sessionConnections method
func TestSessRegSessionConnections(t *testing.T) {
	expectedSessionsNum := 1
	reg := newSessionRegistry(0, nil)

	// Register first connection on session A
	cltA1 := newConnection(nil, []byte("A1"), nil, ConnectionOptions{})
	sessA1 := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA1

	require.NoError(t, reg.register(cltA1))

	// Register second connection on same session
	cltA2 := newConnection(nil, []byte("A2"), nil, ConnectionOptions{})
	sessA2 := NewSession(nil, func() string { return "testkey_A" })
	cltA2.session = &sessA2

	require.NoError(t, reg.register(cltA2))

	// Expect 1 active session
	require.Equal(t, expectedSessionsNum, reg.activeSessionsNum())

	// Expect 2 active connections on session A
	list := reg.sessionConnections("testkey_A")
	require.Len(t, list, 2)
	require.Contains(t, list, cltA1)
	require.Contains(t, list, cltA2)
}

// TestSessRegDestruction tests deregistration
func TestSessRegDestruction(t *testing.T) {
	cb := [2]bool{false, false}
	reg := newSessionRegistry(0, func(sessionKey string) {
		if sessionKey == "testkey_A" {
			cb[0] = true
		} else if sessionKey == "testkey_B" {
			cb[1] = true
		}
	})

	// Register 2 connections on session A, and 1 on session B
	cltA1 := newConnection(nil, nil, nil, ConnectionOptions{})
	sessA := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA

	cltA2 := newConnection(nil, nil, nil, ConnectionOptions{})
	cltA2.session = &sessA

	cltB1 := newConnection(nil, nil, nil, ConnectionOptions{})
	sessB := NewSession(nil, func() string { return "testkey_B" })
	cltB1.session = &sessB

	require.NoError(t, reg.register(cltA1))
	require.NoError(t, reg.register(cltA2))
	require.NoError(t, reg.register(cltB1))

	// Deregister both
	reg.deregister(cltA1, true)
	require.Equal(t, [2]bool{false, false}, cb)

	reg.deregister(cltA2, true)
	require.Equal(t, [2]bool{true, false}, cb)

	reg.deregister(cltB1, true)
	require.Equal(t, [2]bool{true, true}, cb)
}
