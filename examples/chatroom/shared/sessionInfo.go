package shared

import (
	webwire "github.com/qbeon/webwire-go"
)

var sessionInfoFieldNames = []string{"username"}

// SessionInfo implements the webwire.SessionInfo interface
// for this particular example
type SessionInfo struct {
	Username string
}

// Copy implements the webwire.SessionInfo interface.
// It deep-copies the object and returns it's exact clone
func (sinf *SessionInfo) Copy() webwire.SessionInfo {
	return &SessionInfo{
		Username: sinf.Username,
	}
}

// Fields implements the webwire.SessionInfo interface.
// It returns a constant list of the names of all fields of the object
func (sinf *SessionInfo) Fields() []string {
	return sessionInfoFieldNames
}

// Value implements the webwire.SessionInfo interface.
// It returns an exact deep copy of a session info field value
func (sinf *SessionInfo) Value(fieldName string) interface{} {
	switch fieldName {
	case "username":
		return sinf.Username
	}
	return nil
}

// SessionInfoParser parses the given session info data into a
// webwire.SessionInfo compliant object specific to this application
func SessionInfoParser(data map[string]interface{}) webwire.SessionInfo {
	return &SessionInfo{
		Username: data["username"].(string),
	}
}
