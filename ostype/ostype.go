package ostype

import (
	"encoding/json"
	"fmt"
	"strings"
)

// OperatingSystem represents the type and version of the clients OS
type OperatingSystem byte

const (
	// Unknown represents an unknown operating system
	Unknown OperatingSystem = iota
	// Windows10 represents Microsoft Windows 10
	Windows10
	// Windows81 represents Microsoft Windows 8.1
	Windows81
	// Windows8 represents Microsoft Windows 8
	Windows8
	// Windows7 represents Microsoft Windows 7
	Windows7
	// WindowsVista represents Microsoft Windows Vista
	WindowsVista
	// WindowsXP represents Microsoft Windows XP
	WindowsXP
	// WindowsNT represents Microsoft Windows NT
	WindowsNT
	// Windows2000 represents Microsoft Windows 2000
	Windows2000
	// MacOSX1013 represents Apple Mac OSX 10.13
	MacOSX1013
	// MacOSX1012 represents Apple Mac OSX 10.12
	MacOSX1012
	// MacOSX1010 represents Apple Mac OSX 10.10
	MacOSX1010
	// MacOSX109 represents Apple Mac OSX 10.9
	MacOSX109
	// MacOSX108 represents Apple Mac OSX 10.8
	MacOSX108
	// MacOSX107 represents Apple Mac OSX 10.7
	MacOSX107
	// MacOSX106 represents Apple Mac OSX 10.6
	MacOSX106
	// MacOSX105 represents Apple Mac OSX 10.5
	MacOSX105
	// Linux represents an unknown Linux distribution
	Linux
	// Android8 represents Google Android 8 Oreo
	Android8
	// Android7 represents Google Android 7 Nougat
	Android7
	// Android6 represents Google Android 6 Marshmallow
	Android6
	// Android5 represents Google Android 5 Lollipop
	Android5
	// Android44 represents Google Android 4.4 Kit Kat
	Android44
	// Android41 represents Google Android 4.1 Jelly Bean
	Android41
	// Android4 represents Google Android 4 Ice Cream Sandwich
	Android4
	// Android23 represents Google Android 2.3 Gingerbread
	Android23
	// IOS11 represents Apple IOS 11 Tigris
	IOS11
	// IOS10 represents Apple IOS 10 Whitetail
	IOS10
	// IOS9 represents Apple IOS 9 Monarch
	IOS9
	// IOS8 represents Apple IOS 8 Okemo
	IOS8
	// IOS7 represents Apple IOS 7 Innsbruck
	IOS7
	// IOS6 represents Apple IOS 6 Sundance
	IOS6
	// IOS5 represents Apple IOS 5 Telluride
	IOS5
	// IOS4 represents Apple IOS 4 Apex
	IOS4
)

var valToStr = map[OperatingSystem]string{
	Unknown:      "unknown",
	Windows10:    "windows10",
	Windows81:    "windows81",
	Windows8:     "windows8",
	Windows7:     "windows7",
	WindowsVista: "windowsvista",
	WindowsXP:    "windowsxp",
	WindowsNT:    "windowsnt",
	Windows2000:  "windows2000",
	MacOSX1013:   "macosx1013",
	MacOSX1012:   "macosx1012",
	MacOSX1010:   "macosx1010",
	MacOSX109:    "macosx109",
	MacOSX108:    "macosx108",
	MacOSX107:    "macosx107",
	MacOSX106:    "macosx106",
	MacOSX105:    "macosx105",
	Linux:        "linux",
	Android8:     "android8",
	Android7:     "android7",
	Android6:     "android6",
	Android5:     "android5",
	Android44:    "android44",
	Android41:    "android41",
	Android4:     "android4",
	Android23:    "android23",
	IOS11:        "ios11",
	IOS10:        "ios10",
	IOS9:         "ios9",
	IOS8:         "ios8",
	IOS7:         "ios7",
	IOS6:         "ios6",
	IOS5:         "ios5",
	IOS4:         "ios4",
}

var strToVal = map[string]OperatingSystem{
	"unknown":      Unknown,
	"windows10":    Windows10,
	"windows81":    Windows81,
	"windows8":     Windows8,
	"windows7":     Windows7,
	"windowsvista": WindowsVista,
	"windowsxp":    WindowsXP,
	"windowsnt":    WindowsNT,
	"windows2000":  Windows2000,
	"macosx1013":   MacOSX1013,
	"macosx1012":   MacOSX1012,
	"macosx1010":   MacOSX1010,
	"macosx109":    MacOSX109,
	"macosx108":    MacOSX108,
	"macosx107":    MacOSX107,
	"macosx106":    MacOSX106,
	"macosx105":    MacOSX105,
	"linux":        Linux,
	"android8":     Android8,
	"android7":     Android7,
	"android6":     Android6,
	"android5":     Android5,
	"android44":    Android44,
	"android41":    Android41,
	"android4":     Android4,
	"android23":    Android23,
	"ios11":        IOS11,
	"ios10":        IOS10,
	"ios9":         IOS9,
	"ios8":         IOS8,
	"ios7":         IOS7,
	"ios6":         IOS6,
	"ios5":         IOS5,
	"ios4":         IOS4,
}

// Turns the operating system identifier into a string
func (os OperatingSystem) String() string {
	return valToStr[os]
}

// FromString initialized
func (os *OperatingSystem) FromString(str string) error {
	str = strings.ToLower(str)
	var exists bool
	if *os, exists = strToVal[str]; !exists {
		return fmt.Errorf("Invalid operating system type "+
			"string representation: %s",
			str,
		)
	}
	return nil
}

// MarshalJSON implements the JSON marshalling interface
func (os OperatingSystem) MarshalJSON() ([]byte, error) {
	return json.Marshal(os.String())
}

// UnmarshalJSON implements the JSON unmarshalling interface
func (os *OperatingSystem) UnmarshalJSON(bytes []byte) error {
	return os.FromString(string(bytes[1 : len(bytes)-1]))
}
