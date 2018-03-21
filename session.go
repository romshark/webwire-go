package webwire

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

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

// GenerateSessionKey returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateSessionKey() string {
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
	Key          string      `json:"key"`
	CreationDate time.Time   `json:"crt"`
	Info         interface{} `json:"inf"`
}

// NewSession generates a new session object
// generating a cryptographically random secure key
func NewSession(info interface{}, customGenerator func() string) Session {
	var key string
	if customGenerator == nil {
		// Use default generator
		key = GenerateSessionKey()
	} else {
		key = customGenerator()
		if len(key) < 1 {
			panic(fmt.Errorf(
				"Invalid session key returned by custom session key generator (empty)",
			))
		}
	}
	return Session{
		key,
		time.Now(),
		info,
	}
}
