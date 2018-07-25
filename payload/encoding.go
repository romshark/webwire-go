package payload

// Encoding represents the type of encoding of the message payload
type Encoding int

const (
	// Binary represents unencoded binary data
	Binary Encoding = iota

	// Utf8 represents UTF8 encoding
	Utf8

	// Utf16 represents UTF16 encoding
	Utf16
)

// String stringifies the encoding type
func (enc Encoding) String() string {
	switch enc {
	case Binary:
		return "binary"
	case Utf8:
		return "utf8"
	case Utf16:
		return "utf16"
	}
	return ""
}
