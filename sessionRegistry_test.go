package webwire

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

// TestSessRegRegisteration tests registeration
func TestSessRegRegisteration(t *testing.T) {
	reg := newSessionRegistry(0)

	// Register connection with session
	clt := newConnection(nil, "", nil)
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
	reg := newSessionRegistry(0)

	// Register 2 connections on two separate sessions
	cltA1 := newConnection(nil, "", nil)
	sessA := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA

	cltB1 := newConnection(nil, "", nil)
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
	reg := newSessionRegistry(0)

	// Register first connection on session A
	cltA1 := newConnection(nil, "", nil)
	sessA1 := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA1

	require.NoError(t, reg.register(cltA1))

	// Register second connection on same session
	cltA2 := newConnection(nil, "", nil)
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
	reg := newSessionRegistry(1)

	// Register first connection on session A
	cltA1 := newConnection(nil, "", nil)
	sessA1 := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA1

	require.NoError(t, reg.register(cltA1))

	// Register first connection on session A
	cltA2 := newConnection(nil, "", nil)
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
	reg := newSessionRegistry(0)

	// Register 2 connections on two separate sessions
	cltA1 := newConnection(nil, "", nil)
	sessA := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA

	cltB1 := newConnection(nil, "", nil)
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
	require.Equal(t, 0, reg.deregister(cltA1))

	// Expect 1 active session after deregistration of the first connection
	require.Equal(t, 1, reg.activeSessionsNum())

	// Expect 0 active connections on deregistered session
	require.Equal(t, -1, reg.sessionConnectionsNum("testkey_A"))
	require.Equal(t, 1, reg.sessionConnectionsNum("testkey_B"))

	// Deregister second connection, expect 0 because session was removed
	require.Equal(t, 0, reg.deregister(cltB1))

	// Expect no active sessions after deregistration of the second connection
	require.Equal(t, 0, reg.activeSessionsNum())

	// Expect no active connections on both sessions
	require.Equal(t, -1, reg.sessionConnectionsNum("testkey_A"))
	require.Equal(t, -1, reg.sessionConnectionsNum("testkey_B"))
}

// TestSessRegDeregistrationMultiple tests deregistration of multiple
// connections of a single session
func TestSessRegDeregistrationMultiple(t *testing.T) {
	reg := newSessionRegistry(0)

	// Register 2 connections on two separate sessions
	cltA1 := newConnection(nil, "", nil)
	sessA1 := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA1

	cltA2 := newConnection(nil, "", nil)
	sessA2 := NewSession(nil, func() string { return "testkey_A" })
	cltA2.session = &sessA2

	require.NoError(t, reg.register(cltA1))
	require.NoError(t, reg.register(cltA2))

	// Expect 1 active session
	require.Equal(t, 1, reg.activeSessionsNum())

	// Expect 2 active connections on session A
	require.Equal(t, 2, reg.sessionConnectionsNum("testkey_A"))

	// Deregister first connection, expect 1
	require.Equal(t, 1, reg.deregister(cltA1))

	// Still expect 1 active session
	// after deregistration of the first connection
	require.Equal(t, 1, reg.activeSessionsNum())

	// Expect 1 remaining active connection on deregistered session
	require.Equal(t, 1, reg.sessionConnectionsNum("testkey_A"))

	// Deregister second connection, expect 0 because session was removed
	require.Equal(t, 0, reg.deregister(cltA2))

	// Expect no active sessions after deregistration of the second connection
	require.Equal(t, 0, reg.activeSessionsNum())

	// Expect no active connections session A
	require.Equal(t, -1, reg.sessionConnectionsNum("testkey_A"))
}

// TestSessRegConcurrentAccess tests concurrent (de)registeration
func TestSessRegConcurrentAccess(t *testing.T) {
	reg := newSessionRegistry(0)
	connsToRegister := uint(16)
	registeredConns := make([]*connection, connsToRegister)
	var awaitRegistration sync.WaitGroup
	var awaitDeregistration sync.WaitGroup

	// Populate registered conns map
	for i := uint(0); i < connsToRegister; i++ {
		// Create a connection on session A
		clt := newConnection(nil, "", nil)
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
			result := reg.deregister(registeredConns[index])
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
	reg := newSessionRegistry(0)

	// Register first connection on session A
	cltA1 := newConnection(nil, "A1", nil)
	sessA1 := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA1

	require.NoError(t, reg.register(cltA1))

	// Register second connection on same session
	cltA2 := newConnection(nil, "A2", nil)
	sessA2 := NewSession(nil, func() string { return "testkey_A" })
	cltA2.session = &sessA2

	require.NoError(t, reg.register(cltA2))

	// Expect 1 active session
	require.Equal(t, expectedSessionsNum, reg.activeSessionsNum())

	// Expect 2 active connections on session A
	list := reg.sessionConnections("testkey_A")
	require.Len(t, list, 2)
	require.Equal(t, cltA1, list[0])
	require.Equal(t, cltA2, list[1])
}
