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
	UNKNOWN OperatingSystem = iota
	WINDOWS_10
	WINDOWS_8_1
	WINDOWS_8
	WINDOWS_7
	WINDOWS_VISTA
	WINDOWS_XP
	WINDOWS_NT
	WINDOWS_2000
	MACOSX_10_13
	MACOSX_10_12
	MACOSX_10_10
	MACOSX_10_9
	MACOSX_10_8
	MACOSX_10_7
	MACOSX_10_6
	MACOSX_10_5
	LINUX
	ANDROID_8
	ANDROID_7
	ANDROID_6
	ANDROID_5
	ANDROID_4_4
	ANDROID_4
	ANDROID_2_3
	IOS_11
	IOS_10
	IOS_9
	IOS_8
	IOS_7
	IOS_6
	IOS_5
	IOS_4
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
