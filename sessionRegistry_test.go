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

	if err := reg.register(clt); err != nil {
		t.Fatal("Unexpected error: %s", err)
	}

	// Expect 1 active session
	if reg.activeSessionsNum() != 1 {
		t.Fatal("Expected ActiveSessionsNum to return 1")
	}

	// Expect 1 active connection on session A
	if reg.sessionConnectionsNum("testkey_A") != 1 {
		t.Fatal("Expected ActiveSessionsNum to return 1")
	}
}

// TestSessRegActiveSessionsNum tests the ActiveSessionsNum method
func TestSessRegActiveSessionsNum(t *testing.T) {
	expectedSessionsNum := 2
	reg := newSessionRegistry(0)

	// Register 2 connections on two separate sessions
	cltA1 := newClientAgent(nil, "", nil)
	sessA := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA

	cltB1 := newClientAgent(nil, "", nil)
	sessB := NewSession(nil, func() string { return "testkey_B" })
	cltB1.session = &sessB

	if err := reg.register(cltA1); err != nil {
		t.Fatal("Unexpected error: %s", err)
	}
	if err := reg.register(cltB1); err != nil {
		t.Fatal("Unexpected error: %s", err)
	}

	// Expect 2 active sessions
	if reg.activeSessionsNum() != expectedSessionsNum {
		t.Fatalf("Expected ActiveSessionsNum to return %d", expectedSessionsNum)
	}

	// Expect 1 active connection on each session
	if reg.sessionConnectionsNum("testkey_A") != 1 {
		t.Fatal("Expected ActiveSessionsNum to return 1")
	}
	if reg.sessionConnectionsNum("testkey_B") != 1 {
		t.Fatal("Expected ActiveSessionsNum to return 1")
	}
}

// TestSessRegsessionConnectionsNum tests the sessionConnectionsNum method
func TestSessRegsessionConnectionsNum(t *testing.T) {
	expectedSessionsNum := 1
	reg := newSessionRegistry(0)

	// Register first connection on session A
	cltA1 := newClientAgent(nil, "", nil)
	sessA1 := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA1

	if err := reg.register(cltA1); err != nil {
		t.Fatal("Unexpected error: %s", err)
	}

	// Register second connection on same session
	cltA2 := newClientAgent(nil, "", nil)
	sessA2 := NewSession(nil, func() string { return "testkey_A" })
	cltA2.session = &sessA2

	if err := reg.register(cltA2); err != nil {
		t.Fatal("Unexpected error: %s", err)
	}

	// Expect 1 active session
	if reg.activeSessionsNum() != expectedSessionsNum {
		t.Fatalf("Expected ActiveSessionsNum to return %d", expectedSessionsNum)
	}

	// Expect 2 active connections on session A
	if reg.sessionConnectionsNum("testkey_A") != 2 {
		t.Fatal("Expected ActiveSessionsNum to return 2")
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

	if err := reg.register(cltA1); err != nil {
		t.Fatal("Unexpected error: %s", err)
	}

	// Register first connection on session A
	cltA2 := newClientAgent(nil, "", nil)
	sessA2 := NewSession(nil, func() string { return "testkey_A" })
	cltA2.session = &sessA2

	if reg.register(cltA1) == nil {
		t.Fatal("Expected register to return an error " +
			"due to the limit of concurrent connection being reached",
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

	if err := reg.register(cltA1); err != nil {
		t.Fatal("Unexpected error: %s", err)
	}
	if err := reg.register(cltB1); err != nil {
		t.Fatal("Unexpected error: %s", err)
	}

	// Expect 2 active sessions
	if reg.activeSessionsNum() != 2 {
		t.Fatal("Expected ActiveSessionsNum to return 2")
	}

	// Expect 1 active connection on each session
	if reg.sessionConnectionsNum("testkey_A") != 1 {
		t.Fatal("Expected ActiveSessionsNum to return 1")
	}
	if reg.sessionConnectionsNum("testkey_B") != 1 {
		t.Fatal("Expected ActiveSessionsNum to return 1")
	}

	// Deregister first connection, expect 0 because session was removed
	if reg.deregister(cltA1) != 0 {
		t.Fatal("Expected deregister to return 0")
	}

	// Expect 1 active session after deregistration of the first connection
	if reg.activeSessionsNum() != 1 {
		t.Fatal("Expected ActiveSessionsNum to return 1")
	}

	// Expect 0 active connections on deregistered session
	if reg.sessionConnectionsNum("testkey_A") != -1 {
		t.Fatal("Expected ActiveSessionsNum to return -1")
	}
	if reg.sessionConnectionsNum("testkey_B") != 1 {
		t.Fatal("Expected ActiveSessionsNum to return 1")
	}

	// Deregister second connection, expect 0 because session was removed
	if reg.deregister(cltB1) != 0 {
		t.Fatal("Expected deregister to return 0")
	}

	// Expect no active sessions after deregistration of the second connection
	if reg.activeSessionsNum() != 0 {
		t.Fatal("Expected ActiveSessionsNum to return 0")
	}

	// Expect no active connections on both sessions
	if reg.sessionConnectionsNum("testkey_A") != -1 {
		t.Fatal("Expected ActiveSessionsNum to return -1")
	}
	if reg.sessionConnectionsNum("testkey_B") != -1 {
		t.Fatal("Expected ActiveSessionsNum to return -1")
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

	if err := reg.register(cltA1); err != nil {
		t.Fatal("Unexpected error: %s", err)
	}
	if err := reg.register(cltA2); err != nil {
		t.Fatal("Unexpected error: %s", err)
	}

	// Expect 1 active session
	if reg.activeSessionsNum() != 1 {
		t.Fatal("Expected ActiveSessionsNum to return 1")
	}

	// Expect 2 active connections on session A
	if reg.sessionConnectionsNum("testkey_A") != 2 {
		t.Fatal("Expected ActiveSessionsNum to return 2")
	}

	// Deregister first connection, expect 1
	if reg.deregister(cltA1) != 1 {
		t.Fatal("Expected deregister to return 1")
	}

	// Still expect 1 active session
	// after deregistration of the first connection
	if reg.activeSessionsNum() != 1 {
		t.Fatal("Expected ActiveSessionsNum to return 1")
	}

	// Expect 1 remaining active connection on deregistered session
	if reg.sessionConnectionsNum("testkey_A") != 1 {
		t.Fatal("Expected ActiveSessionsNum to return 0")
	}

	// Deregister second connection, expect 0 because session was removed
	if reg.deregister(cltA2) != 0 {
		t.Fatal("Expected deregister to return 0")
	}

	// Expect no active sessions after deregistration of the second connection
	if reg.activeSessionsNum() != 0 {
		t.Fatal("Expected ActiveSessionsNum to return 0")
	}

	// Expect no active connections session A
	if reg.sessionConnectionsNum("testkey_A") != -1 {
		t.Fatal("Expected ActiveSessionsNum to return -1")
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
			if err := reg.register(registeredConns[index]); err != nil {
				t.Error("Unexpected error: %s", err)
			}
			awaitRegistration.Done()
		}()
	}

	// Wait for all goroutines to finish before evaluating the results
	awaitRegistration.Wait()

	// Expect 1 active session
	if reg.activeSessionsNum() != 1 {
		t.Fatal("Expected ActiveSessionsNum to return 1")
	}

	// Expect N active connections on session A
	actualSessCons := reg.sessionConnectionsNum("testkey_A")
	if actualSessCons < 0 {
		t.Fatalf("Missing session testkey_A")
	}
	if uint(actualSessCons) != connsToRegister {
		t.Fatalf(
			"Expected ActiveSessionsNum to return %d, got: %d",
			connsToRegister,
			actualSessCons,
		)
	}

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

	if negativeReturnCounter > 1 {
		t.Fatalf(
			"deregister returned false %d times, expected 1",
			negativeReturnCounter,
		)
	}

	// Expect 1 active session
	if reg.activeSessionsNum() != 0 {
		t.Fatal("Expected ActiveSessionsNum to return 0")
	}

	// Expect no active connections on session A
	actualSessConsAfter := reg.sessionConnectionsNum("testkey_A")
	if actualSessConsAfter != -1 {
		t.Fatalf(
			"Expected ActiveSessionsNum to return 0, got: %d",
			actualSessConsAfter,
		)
	}
}

// TestSessRegSessionConnections tests the sessionConnections method
func TestSessRegSessionConnections(t *testing.T) {
	expectedSessionsNum := 1
	reg := newSessionRegistry(0)

	// Register first connection on session A
	cltA1 := newClientAgent(nil, "A1", nil)
	sessA1 := NewSession(nil, func() string { return "testkey_A" })
	cltA1.session = &sessA1

	if err := reg.register(cltA1); err != nil {
		t.Fatal("Unexpected error: %s", err)
	}

	// Register second connection on same session
	cltA2 := newClientAgent(nil, "A2", nil)
	sessA2 := NewSession(nil, func() string { return "testkey_A" })
	cltA2.session = &sessA2

	if err := reg.register(cltA2); err != nil {
		t.Fatal("Unexpected error: %s", err)
	}

	// Expect 1 active session
	if reg.activeSessionsNum() != expectedSessionsNum {
		t.Fatalf("Expected ActiveSessionsNum to return %d", expectedSessionsNum)
	}

	// Expect 2 active connections on session A
	list := reg.sessionConnections("testkey_A")
	if len(list) != 2 {
		t.Fatal("Expected 2 active connections on session testkey_A")
	}
	if list[0] != cltA1 {
		t.Fatal("Expected cltA1 to be in the list of active connections")
	}
	if list[1] != cltA2 {
		t.Fatal("Expected cltA2 to be in the list of active connections")
	}
}
