package webwire

import (
	"fmt"
	"encoding/json"
)

// OperatingSystem represents the type and version of the clients OS
type OperatingSystem int
const (
	// OsUnknown represents an unknown operating system
	OsUnknown OperatingSystem = iota
	// OsWindows10 represents Microsoft Windows 10
	OsWindows10
	// OsWindows81 represents Microsoft Windows 8.1
	OsWindows81
	// OsWindows8 represents Microsoft Windows 8
	OsWindows8
	// OsWindows7 represents Microsoft Windows 7
	OsWindows7
	// OsWindowsVista represents Microsoft Windows Vista
	OsWindowsVista
	// OsWindowsXP represents Microsoft Windows XP
	OsWindowsXP
	// OsWindowsNT represents Microsoft Windows NT
	OsWindowsNT
	// OsWindows2000 represents Microsoft Windows 2000
	OsWindows2000
	// OsMacOSX1013 represents Apple Mac OSX 10.13
	OsMacOSX1013
	// OsMacOSX1012 represents Apple Mac OSX 10.12
	OsMacOSX1012
	// OsMacOSX1010 represents Apple Mac OSX 10.10
	OsMacOSX1010
	// OsMacOSX109 represents Apple Mac OSX 10.9
	OsMacOSX109
	// OsMacOSX108 represents Apple Mac OSX 10.8
	OsMacOSX108
	// OsMacOSX107 represents Apple Mac OSX 10.7
	OsMacOSX107
	// OsMacOSX106 represents Apple Mac OSX 10.6
	OsMacOSX106
	// OsMacOSX105 represents Apple Mac OSX 10.5
	OsMacOSX105
	// OsLinux represents an unknown Linux distribution
	OsLinux
	// OsAndroid8 represents Google Android 8 Oreo
	OsAndroid8
	// OsAndroid7 represents Google Android 7 Nougat
	OsAndroid7
	// OsAndroid6 represents Google Android 6 Marshmallow
	OsAndroid6
	// OsAndroid5 represents Google Android 5 Lollipop
	OsAndroid5
	// OsAndroid44 represents Google Android 4.4 Kit Kat
	OsAndroid44
	// OsAndroid41 represents Google Android 4.1 Jelly Bean
	OsAndroid41
	// OsAndroid4 represents Google Android 4 Ice Cream Sandwich
	OsAndroid4
	// OsAndroid23 represents Google Android 2.3 Gingerbread
	OsAndroid23
	// OsIOS11 represents Apple IOS 11 Tigris
	OsIOS11
	// OsIOS10 represents Apple IOS 10 Whitetail
	OsIOS10
	// OsIOS9 represents Apple IOS 9 Monarch
	OsIOS9
	// OsIOS8 represents Apple IOS 8 Okemo
	OsIOS8
	// OsIOS7 represents Apple IOS 7 Innsbruck
	OsIOS7
	// OsIOS6 represents Apple IOS 6 Sundance
	OsIOS6
	// OsIOS5 represents Apple IOS 5 Telluride
	OsIOS5
	// OsIOS4 represents Apple IOS 4 Apex
	OsIOS4
)

func (os OperatingSystem) String() string {
	switch(os) {
	case OsUnknown: return "OsUnknown"
	case OsWindows10: return "OsWindows10"
	case OsWindows81: return "OsWindows81"
	case OsWindows8: return "OsWindows8"
	case OsWindows7: return "OsWindows7"
	case OsWindowsVista: return "OsWindowsVista"
	case OsWindowsXP: return "OsWindowsXP"
	case OsWindowsNT: return "OsWindowsNT"
	case OsWindows2000: return "OsWindows2000"
	case OsMacOSX1013: return "OsMacOSX1013"
	case OsMacOSX1012: return "OsMacOSX1012"
	case OsMacOSX1010: return "OsMacOSX1010"
	case OsMacOSX109: return "OsMacOSX109"
	case OsMacOSX108: return "OsMacOSX108"
	case OsMacOSX107: return "OsMacOSX107"
	case OsMacOSX106: return "OsMacOSX106"
	case OsMacOSX105: return "OsMacOSX105"
	case OsLinux: return "OsLinux"
	case OsAndroid8: return "OsAndroid8"
	case OsAndroid7: return "OsAndroid7"
	case OsAndroid6: return "OsAndroid6"
	case OsAndroid5: return "OsAndroid5"
	case OsAndroid44: return "OsAndroid44"
	case OsAndroid41: return "OsAndroid41"
	case OsAndroid4: return "OsAndroid4"
	case OsAndroid23: return "OsAndroid23"
	case OsIOS11: return "OsIOS11"
	case OsIOS10: return "OsIOS10"
	case OsIOS9: return "OsIOS9"
	case OsIOS8: return "OsIOS8"
	case OsIOS7: return "OsIOS7"
	case OsIOS6: return "OsIOS6"
	case OsIOS5: return "OsIOS5"
	case OsIOS4: return "OsIOS4"
	default: panic(fmt.Errorf("Invalid OperatingSystem numeric value: %d", os))
	}
}

// FromString tries to parse the operating system type from the given string
func (os *OperatingSystem) FromString(str string) error {
	switch(str) {
	case "OsUnknown": *os = OsUnknown
	case "OsWindows10": *os = OsWindows10
	case "OsWindows81": *os = OsWindows81
	case "OsWindows8": *os = OsWindows8
	case "OsWindows7": *os = OsWindows7
	case "OsWindowsVista": *os = OsWindowsVista
	case "OsWindowsXP": *os = OsWindowsXP
	case "OsWindowsNT": *os = OsWindowsNT
	case "OsWindows2000": *os = OsWindows2000
	case "OsMacOSX1013": *os = OsMacOSX1013
	case "OsMacOSX1012": *os = OsMacOSX1012
	case "OsMacOSX1010": *os = OsMacOSX1010
	case "OsMacOSX109": *os = OsMacOSX109
	case "OsMacOSX108": *os = OsMacOSX108
	case "OsMacOSX107": *os = OsMacOSX107
	case "OsMacOSX106": *os = OsMacOSX106
	case "OsMacOSX105": *os = OsMacOSX105
	case "OsLinux": *os = OsLinux
	case "OsAndroid8": *os = OsAndroid8
	case "OsAndroid7": *os = OsAndroid7
	case "OsAndroid6": *os = OsAndroid6
	case "OsAndroid5": *os = OsAndroid5
	case "OsAndroid44": *os = OsAndroid44
	case "OsAndroid41": *os = OsAndroid41
	case "OsAndroid4": *os = OsAndroid4
	case "OsAndroid23": *os = OsAndroid23
	case "OsIOS11": *os = OsIOS11
	case "OsIOS10": *os = OsIOS10
	case "OsIOS9": *os = OsIOS9
	case "OsIOS8": *os = OsIOS8
	case "OsIOS7": *os = OsIOS7
	case "OsIOS6": *os = OsIOS6
	case "OsIOS5": *os = OsIOS5
	case "OsIOS4": *os = OsIOS4
	default: return fmt.Errorf("Invalid OperatingSystem string value: %d", os)
	}
	return nil
}

// MarshalJSON implements the JSON marshalling interface
func (os OperatingSystem) MarshalJSON() ([]byte, error) {
	return json.Marshal(os.String())
}

// UnmarshalJSON implements the JSON unmarshalling interface
func (os *OperatingSystem) UnmarshalJSON(bytes []byte) error {
	return os.FromString(string(bytes[1 : len(bytes) - 1]))
}
