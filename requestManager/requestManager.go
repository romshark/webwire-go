package requestmanager

import (
	"encoding/binary"
	"sync"
	"sync/atomic"

	"github.com/qbeon/webwire-go/message"
)

// RequestManager manages and keeps track of outgoing pending requests
type RequestManager struct {
	lastID uint64
	lock   *sync.RWMutex

	// pending represents an indexed list of all pending requests
	pending map[[8]byte]*Request
}

// NewRequestManager constructs and returns a new instance of a RequestManager
func NewRequestManager() RequestManager {
	return RequestManager{
		lastID:  0,
		lock:    &sync.RWMutex{},
		pending: make(map[[8]byte]*Request),
	}
}

// Create creates and registers a new request.
// Create doesn't start the timeout timer,
// this is done in the subsequent request.AwaitReply
func (manager *RequestManager) Create() *Request {
	// Generate unique request identifier by incrementing the last assigned id
	ident := atomic.AddUint64(&manager.lastID, 1)

	identBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(identBytes, ident)
	newRequest := &Request{
		manager:         manager,
		IdentifierBytes: identBytes,
		Reply:           make(chan genericReply, 1),
	}
	copy(newRequest.Identifier[:], identBytes)

	// Register the newly created request
	manager.lock.Lock()
	manager.pending[newRequest.Identifier] = newRequest
	manager.lock.Unlock()

	return newRequest
}

// deregister deregisters the given clients session from the list
// of currently pending requests
func (manager *RequestManager) deregister(identifier [8]byte) {
	manager.lock.Lock()
	delete(manager.pending, identifier)
	manager.lock.Unlock()
}

// Fulfill fulfills the request associated with the given request identifier
// with the provided reply payload.
// Returns true if a pending request was fulfilled and deregistered,
// otherwise returns false
func (manager *RequestManager) Fulfill(msg *message.Message) bool {
	manager.lock.RLock()
	req, exists := manager.pending[msg.MsgIdentifier]
	manager.lock.RUnlock()

	if !exists {
		return false
	}

	manager.deregister(msg.MsgIdentifier)
	req.Reply <- genericReply{
		ReplyMsg: msg,
	}
	return true
}

// Fail fails the request associated with the given request identifier
// with the provided error. Returns true if a pending request
// was failed and deregistered, otherwise returns false
func (manager *RequestManager) Fail(
	identifier [8]byte,
	err error,
) bool {
	manager.lock.RLock()
	req, exists := manager.pending[identifier]
	manager.lock.RUnlock()

	if !exists {
		return false
	}

	manager.deregister(identifier)
	req.Reply <- genericReply{
		Error: err,
	}
	return true
}

// PendingRequests returns the number of currently pending requests
func (manager *RequestManager) PendingRequests() int {
	manager.lock.RLock()
	len := len(manager.pending)
	manager.lock.RUnlock()
	return len
}

// IsPending returns true if the request associated
// with the given identifier is pending
func (manager *RequestManager) IsPending(identifier [8]byte) bool {
	manager.lock.RLock()
	_, exists := manager.pending[identifier]
	manager.lock.RUnlock()
	return exists
}
