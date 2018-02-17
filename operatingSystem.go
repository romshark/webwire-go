package webwire

import (
	"fmt"
	"encoding/json"
)

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

func (os OperatingSystem) String() string {
	switch(os) {
	case Os_UNKNOWN: return "Os_UNKNOWN"
	case Os_WINDOWS_10: return "Os_WINDOWS_10"
	case Os_WINDOWS_8_1: return "Os_WINDOWS_8_1"
	case Os_WINDOWS_8: return "Os_WINDOWS_8"
	case Os_WINDOWS_7: return "Os_WINDOWS_7"
	case Os_WINDOWS_VISTA: return "Os_WINDOWS_VISTA"
	case Os_WINDOWS_XP: return "Os_WINDOWS_XP"
	case Os_WINDOWS_NT: return "Os_WINDOWS_NT"
	case Os_WINDOWS_2000: return "Os_WINDOWS_2000"
	case Os_MACOSX_10_13: return "Os_MACOSX_10_13"
	case Os_MACOSX_10_12: return "Os_MACOSX_10_12"
	case Os_MACOSX_10_10: return "Os_MACOSX_10_10"
	case Os_MACOSX_10_9: return "Os_MACOSX_10_9"
	case Os_MACOSX_10_8: return "Os_MACOSX_10_8"
	case Os_MACOSX_10_7: return "Os_MACOSX_10_7"
	case Os_MACOSX_10_6: return "Os_MACOSX_10_6"
	case Os_MACOSX_10_5: return "Os_MACOSX_10_5"
	case Os_LINUX: return "Os_LINUX"
	case Os_ANDROID_8: return "Os_ANDROID_8"
	case Os_ANDROID_7: return "Os_ANDROID_7"
	case Os_ANDROID_6: return "Os_ANDROID_6"
	case Os_ANDROID_5: return "Os_ANDROID_5"
	case Os_ANDROID_4_4: return "Os_ANDROID_4_4"
	case Os_ANDROID_4: return "Os_ANDROID_4"
	case Os_ANDROID_2_3: return "Os_ANDROID_2_3"
	case Os_IOS_11: return "Os_IOS_11"
	case Os_IOS_10: return "Os_IOS_10"
	case Os_IOS_9: return "Os_IOS_9"
	case Os_IOS_8: return "Os_IOS_8"
	case Os_IOS_7: return "Os_IOS_7"
	case Os_IOS_6: return "Os_IOS_6"
	case Os_IOS_5: return "Os_IOS_5"
	case Os_IOS_4: return "Os_IOS_4"
	default: panic(fmt.Errorf("Invalid OperatingSystem numeric value: %d", os))
	}
}

func (os *OperatingSystem) FromString(str string) error {
	switch(str) {
	case "Os_UNKNOWN": *os = Os_UNKNOWN
	case "Os_WINDOWS_10": *os = Os_WINDOWS_10
	case "Os_WINDOWS_8_1": *os = Os_WINDOWS_8_1
	case "Os_WINDOWS_8": *os = Os_WINDOWS_8
	case "Os_WINDOWS_7": *os = Os_WINDOWS_7
	case "Os_WINDOWS_VISTA": *os = Os_WINDOWS_VISTA
	case "Os_WINDOWS_XP": *os = Os_WINDOWS_XP
	case "Os_WINDOWS_NT": *os = Os_WINDOWS_NT
	case "Os_WINDOWS_2000": *os = Os_WINDOWS_2000
	case "Os_MACOSX_10_13": *os = Os_MACOSX_10_13
	case "Os_MACOSX_10_12": *os = Os_MACOSX_10_12
	case "Os_MACOSX_10_10": *os = Os_MACOSX_10_10
	case "Os_MACOSX_10_9": *os = Os_MACOSX_10_9
	case "Os_MACOSX_10_8": *os = Os_MACOSX_10_8
	case "Os_MACOSX_10_7": *os = Os_MACOSX_10_7
	case "Os_MACOSX_10_6": *os = Os_MACOSX_10_6
	case "Os_MACOSX_10_5": *os = Os_MACOSX_10_5
	case "Os_LINUX": *os = Os_LINUX
	case "Os_ANDROID_8": *os = Os_ANDROID_8
	case "Os_ANDROID_7": *os = Os_ANDROID_7
	case "Os_ANDROID_6": *os = Os_ANDROID_6
	case "Os_ANDROID_5": *os = Os_ANDROID_5
	case "Os_ANDROID_4_4": *os = Os_ANDROID_4_4
	case "Os_ANDROID_4": *os = Os_ANDROID_4
	case "Os_ANDROID_2_3": *os = Os_ANDROID_2_3
	case "Os_IOS_11": *os = Os_IOS_11
	case "Os_IOS_10": *os = Os_IOS_10
	case "Os_IOS_9": *os = Os_IOS_9
	case "Os_IOS_8": *os = Os_IOS_8
	case "Os_IOS_7": *os = Os_IOS_7
	case "Os_IOS_6": *os = Os_IOS_6
	case "Os_IOS_5": *os = Os_IOS_5
	case "Os_IOS_4": *os = Os_IOS_4
	default: return fmt.Errorf("Invalid OperatingSystem string value: %d", os)
	}
	return nil
}

func (os OperatingSystem) MarshalJSON() ([]byte, error) {
	return json.Marshal(os.String())
}

func (os *OperatingSystem) UnmarshalJSON(bytes []byte) error {
	return os.FromString(string(bytes[1 : len(bytes) - 1]))
}
