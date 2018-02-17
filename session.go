package webwire

import (
	"fmt"
	"time"
	"encoding/base64"
	cryptoRand "crypto/rand"
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

type OperatingSystem int
const (
	Os_UNKNOWN OperatingSystem = iota
	Os_WINDOWS_10
	Os_WINDOWS_8_1
	Os_WINDOWS_8
	Os_WINDOWS_7
	Os_WINDOWS_VISTA
	Os_WINDOWS_XP
	Os_WINDOWS_NT
	Os_WINDOWS_2000
	Os_MACOSX_10_13
	Os_MACOSX_10_12
	Os_MACOSX_10_10
	Os_MACOSX_10_9
	Os_MACOSX_10_8
	Os_MACOSX_10_7
	Os_MACOSX_10_6
	Os_MACOSX_10_5
	Os_LINUX
	Os_ANDROID_8
	Os_ANDROID_7
	Os_ANDROID_6
	Os_ANDROID_5
	Os_ANDROID_4_4
	Os_ANDROID_4
	Os_ANDROID_2_3
	Os_IOS_11
	Os_IOS_10
	Os_IOS_9
	Os_IOS_8
	Os_IOS_7
	Os_IOS_6
	Os_IOS_5
	Os_IOS_4
)

type Session struct {
	Key string
	OperatingSystem OperatingSystem
	UserAgent string
	CreationDate time.Time
	Info interface {}
}

func NewSession(
	operatingSystem OperatingSystem,
	userAgent string,
	info interface {},
) Session {
	return Session {
		GenerateSessionKey(),
		operatingSystem,
		userAgent,
		time.Now(),
		info,
	}
}
