package webwire

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

// SessionInfo represents a map of session info fields of arbitrary types
type SessionInfo = map[string]interface{}

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(length uint32) (bytes []byte, err error) {
	bytes = make([]byte, length)
	_, err = cryptoRand.Read(bytes)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// generateSessionKey returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateSessionKey() string {
	bytes, err := generateRandomBytes(48)
	if err != nil {
		panic(fmt.Errorf("Could not generate a session key"))
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// Session represents a session object.
// If the key is empty the session is invalid.
// Info can contain arbitrary attached data
type Session struct {
	Key      string      `json:"key"`
	Creation time.Time   `json:"crt"`
	Info     SessionInfo `json:"inf"`
}

// NewSession generates a new session object
// generating a cryptographically random secure key
func NewSession(info SessionInfo, generator func() string) Session {
	key := generator()
	if len(key) < 1 {
		panic(fmt.Errorf("Invalid session key returned by the session key generator (empty)"))
	}
	return Session{
		key,
		time.Now(),
		info,
	}
}

// DefaultSessionKeyGenerator implements the webwire.SessionKeyGenerator interface
type DefaultSessionKeyGenerator struct{}

// NewDefaultSessionKeyGenerator constructs a new default session key generator implementation
func NewDefaultSessionKeyGenerator() SessionKeyGenerator {
	return &DefaultSessionKeyGenerator{}
}

// Generate implements the webwire.Sessio
func (gen *DefaultSessionKeyGenerator) Generate() string {
	return generateSessionKey()
}
