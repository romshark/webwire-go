package webwire

import (
	"sync"
	"sync/atomic"
	"testing"
)

// TestSessRegRegisteration tests registeration
func TestSessRegRegisteration(t *testing.T) {
	reg := newSessionRegistry(0)

	// Register connection with session
	clt := newClientAgent(nil, "", nil)
	sess := NewSession(nil, func() string { return "testkey_A" })
	clt.session = &sess

	if !reg.register(clt) {
		t.Fatal("Expected register to return true")
	}

	// Expect 1 active session
	if reg.ActiveSessions() != 1 {
		t.Fatal("Expected ActiveSessions to return 1")
	}

	// Expect 1 active connection on session A
	if reg.SessionConnections("testkey_A") != 1 {
		t.Fatal("Expected ActiveSessions to return 1")
	}
}

// TestSessRegActiveSessions tests the ActiveSessions method
func TestSessRegActiveSessions(t *testing.T) {
	expectedSessionsNum := 2
	reg := newSessionRegistry(0)

	// Register 2 connections on two separate sessions
	cltA1 := newClientAgent(nil, "", nil)
	sessA := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA

	cltB1 := newClientAgent(nil, "", nil)
	sessB := NewSession(nil, func() string { return "testkey_B" })
	cltB1.session = &sessB

	if !reg.register(cltA1) {
		t.Fatal("Expected register to return true")
	}
	if !reg.register(cltB1) {
		t.Fatal("Expected register to return true")
	}

	// Expect 2 active sessions
	if reg.ActiveSessions() != expectedSessionsNum {
		t.Fatalf("Expected ActiveSessions to return %d", expectedSessionsNum)
	}

	// Expect 1 active connection on each session
	if reg.SessionConnections("testkey_A") != 1 {
		t.Fatal("Expected ActiveSessions to return 1")
	}
	if reg.SessionConnections("testkey_B") != 1 {
		t.Fatal("Expected ActiveSessions to return 1")
	}
}

// TestSessRegSessionConnections tests the SessionConnections method
func TestSessRegSessionConnections(t *testing.T) {
	expectedSessionsNum := 1
	reg := newSessionRegistry(0)

	// Register first connection on session A
	cltA1 := newClientAgent(nil, "", nil)
	sessA1 := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA1

	if !reg.register(cltA1) {
		t.Fatal("Expected register to return true")
	}

	// Register second connection on same session
	cltA2 := newClientAgent(nil, "", nil)
	sessA2 := NewSession(nil, func() string { return "testkey_A" })
	cltA2.session = &sessA2

	if !reg.register(cltA2) {
		t.Fatal("Expected register to return true")
	}

	// Expect 1 active session
	if reg.ActiveSessions() != expectedSessionsNum {
		t.Fatalf("Expected ActiveSessions to return %d", expectedSessionsNum)
	}

	// Expect 2 active connections on session A
	if reg.SessionConnections("testkey_A") != 2 {
		t.Fatal("Expected ActiveSessions to return 2")
	}
}

// TestSessRegSessionMaxConns tests the Register method
// when the maximum number of concurrent connections of a session was reached
func TestSessRegSessionMaxConns(t *testing.T) {
	// Set the maximum number of concurrent session connection to 1
	reg := newSessionRegistry(1)

	// Register first connection on session A
	cltA1 := newClientAgent(nil, "", nil)
	sessA1 := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA1

	if !reg.register(cltA1) {
		t.Fatal("Expected register to return true")
	}

	// Register first connection on session A
	cltA2 := newClientAgent(nil, "", nil)
	sessA2 := NewSession(nil, func() string { return "testkey_A" })
	cltA2.session = &sessA2

	if reg.register(cltA1) {
		t.Fatal(
			"Expected register to return false due to the limit being reached",
		)
	}
}

// TestSessRegDeregistration tests deregistration
func TestSessRegDeregistration(t *testing.T) {
	reg := newSessionRegistry(0)

	// Register 2 connections on two separate sessions
	cltA1 := newClientAgent(nil, "", nil)
	sessA := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA

	cltB1 := newClientAgent(nil, "", nil)
	sessB := NewSession(nil, func() string { return "testkey_B" })
	cltB1.session = &sessB

	if !reg.register(cltA1) {
		t.Fatal("Expected register to return true")
	}
	if !reg.register(cltB1) {
		t.Fatal("Expected register to return true")
	}

	// Expect 2 active sessions
	if reg.ActiveSessions() != 2 {
		t.Fatal("Expected ActiveSessions to return 2")
	}

	// Expect 1 active connection on each session
	if reg.SessionConnections("testkey_A") != 1 {
		t.Fatal("Expected ActiveSessions to return 1")
	}
	if reg.SessionConnections("testkey_B") != 1 {
		t.Fatal("Expected ActiveSessions to return 1")
	}

	// Deregister first connection, expect false because session was removed
	if reg.deregister(cltA1) {
		t.Fatal("Expected deregister to return false")
	}

	// Expect 1 active session after deregistration of the first connection
	if reg.ActiveSessions() != 1 {
		t.Fatal("Expected ActiveSessions to return 1")
	}

	// Expect 0 active connections on deregistered session
	if reg.SessionConnections("testkey_A") != 0 {
		t.Fatal("Expected ActiveSessions to return 0")
	}
	if reg.SessionConnections("testkey_B") != 1 {
		t.Fatal("Expected ActiveSessions to return 1")
	}

	// Deregister second connection, expect false because session was removed
	if reg.deregister(cltB1) {
		t.Fatal("Expected deregister to return false")
	}

	// Expect no active sessions after deregistration of the second connection
	if reg.ActiveSessions() != 0 {
		t.Fatal("Expected ActiveSessions to return 0")
	}

	// Expect no active connections on both sessions
	if reg.SessionConnections("testkey_A") != 0 {
		t.Fatal("Expected ActiveSessions to return 0")
	}
	if reg.SessionConnections("testkey_B") != 0 {
		t.Fatal("Expected ActiveSessions to return 0")
	}
}

// TestSessRegDeregistrationMultiple tests deregistration of multiple
// connections of a single session
func TestSessRegDeregistrationMultiple(t *testing.T) {
	reg := newSessionRegistry(0)

	// Register 2 connections on two separate sessions
	cltA1 := newClientAgent(nil, "", nil)
	sessA1 := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA1

	cltA2 := newClientAgent(nil, "", nil)
	sessA2 := NewSession(nil, func() string { return "testkey_A" })
	cltA2.session = &sessA2

	if !reg.register(cltA1) {
		t.Fatal("Expected register to return true")
	}
	if !reg.register(cltA2) {
		t.Fatal("Expected register to return true")
	}

	// Expect 1 active session
	if reg.ActiveSessions() != 1 {
		t.Fatal("Expected ActiveSessions to return 1")
	}

	// Expect 2 active connections on session A
	if reg.SessionConnections("testkey_A") != 2 {
		t.Fatal("Expected ActiveSessions to return 2")
	}

	// Deregister first connection
	if !reg.deregister(cltA1) {
		t.Fatal("Expected deregister to return true")
	}

	// Still expect 1 active session
	// after deregistration of the first connection
	if reg.ActiveSessions() != 1 {
		t.Fatal("Expected ActiveSessions to return 1")
	}

	// Expect 1 remaining active connection on deregistered session
	if reg.SessionConnections("testkey_A") != 1 {
		t.Fatal("Expected ActiveSessions to return 0")
	}

	// Deregister second connection, expect false because session was removed
	if reg.deregister(cltA2) {
		t.Fatal("Expected deregister to return false")
	}

	// Expect no active sessions after deregistration of the second connection
	if reg.ActiveSessions() != 0 {
		t.Fatal("Expected ActiveSessions to return 0")
	}

	// Expect no active connections session A
	if reg.SessionConnections("testkey_A") != 0 {
		t.Fatal("Expected ActiveSessions to return 0")
	}
}

// TestSessRegConcurrentAccess tests concurrent (de)registeration
func TestSessRegConcurrentAccess(t *testing.T) {
	reg := newSessionRegistry(0)
	connsToRegister := uint(16)
	registeredConns := make([]*Client, connsToRegister)
	var awaitRegistration sync.WaitGroup
	var awaitDeregistration sync.WaitGroup

	// Populate registered conns map
	for i := uint(0); i < connsToRegister; i++ {
		// Create a connection on session A
		clt := newClientAgent(nil, "", nil)
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
			if !reg.register(registeredConns[index]) {
				t.Error("Expected register to return true")
			}
			awaitRegistration.Done()
		}()
	}

	// Wait for all goroutines to finish before evaluating the results
	awaitRegistration.Wait()

	// Expect 1 active session
	if reg.ActiveSessions() != 1 {
		t.Fatal("Expected ActiveSessions to return 1")
	}

	// Expect N active connections on session A
	actualSessCons := reg.SessionConnections("testkey_A")
	if actualSessCons != connsToRegister {
		t.Fatalf(
			"Expected ActiveSessions to return %d, got: %d",
			connsToRegister,
			actualSessCons,
		)
	}

	// Concurrently deregister multiple connections from different goroutines
	falseReturnCounter := uint32(0)
	for i := uint(0); i < connsToRegister; i++ {
		index := i
		go func() {
			// Deregister one of the connections
			result := reg.deregister(registeredConns[index])
			if !result {
				atomic.AddUint32(&falseReturnCounter, 1)
			}
			awaitDeregistration.Done()
		}()
	}

	// Wait for all goroutines to finish before evaluating the results
	awaitDeregistration.Wait()

	if falseReturnCounter > 1 {
		t.Fatalf(
			"deregister returned false %d times, expected 1",
			falseReturnCounter,
		)
	}

	// Expect 1 active session
	if reg.ActiveSessions() != 0 {
		t.Fatal("Expected ActiveSessions to return 0")
	}

	// Expect no active connections on session A
	actualSessConsAfter := reg.SessionConnections("testkey_A")
	if actualSessConsAfter != 0 {
		t.Fatalf(
			"Expected ActiveSessions to return 0, got: %d",
			actualSessConsAfter,
		)
	}
}
