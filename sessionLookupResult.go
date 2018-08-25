package webwire

import "time"

// NewSessionLookupResult creates a new result of a session lookup operation
func NewSessionLookupResult(
	creation time.Time,
	lastLookup time.Time,
	info map[string]interface{},
) SessionLookupResult {
	return &sessionLookupResult{
		creation:   creation,
		lastLookup: lastLookup,
		info:       info,
	}
}

// sessionLookupResult represents an implementation
// of the SessionLookupResult interface
type sessionLookupResult struct {
	creation   time.Time
	lastLookup time.Time
	info       map[string]interface{}
}

// Creation implements the SessionLookupResult interface
func (slr *sessionLookupResult) Creation() time.Time {
	return slr.creation
}

// LastLookup implements the SessionLookupResult interface
func (slr *sessionLookupResult) LastLookup() time.Time {
	return slr.lastLookup
}

// Info implements the SessionLookupResult interface
func (slr *sessionLookupResult) Info() map[string]interface{} {
	return slr.info
}
