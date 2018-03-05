package webwire

// PayloadEncoding represents the type of encoding of the message payload
type PayloadEncoding int

const (
	// EncodingBinary represents unencoded binary data
	EncodingBinary PayloadEncoding = iota

	// EncodingUtf8 represents UTF8 encoding
	EncodingUtf8

	// EncodingUtf16 represents UTF16 encoding
	EncodingUtf16
)

func (enc PayloadEncoding) String() string {
	switch enc {
	case EncodingBinary:
		return "binary"
	case EncodingUtf8:
		return "utf8"
	case EncodingUtf16:
		return "utf16"
	}
	return ""
}
