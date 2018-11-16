package payload

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestEncodingStringification tests the stringification method
// of the Encoding enumeration type
func TestEncodingStringification(t *testing.T) {
	binaryEncoding := Binary
	require.Equal(t, "binary", binaryEncoding.String())

	utf8Encoding := Utf8
	require.Equal(t, "utf8", utf8Encoding.String())

	utf16Encoding := Utf16
	require.Equal(t, "utf16", utf16Encoding.String())
}

// TestConvertUtf8ToUtf8 tests the Utf8() payload conversion method
// with a payload already encoded in UTF8
func TestConvertUtf8ToUtf8(t *testing.T) {
	payload := Payload{
		Encoding: Utf8,
		Data:     []byte{65, 66, 67}, // "ABC"
	}

	result, err := payload.Utf8()
	require.NoError(t, err)
	require.Equal(t, []byte("ABC"), result)
}

// TestConvertBinaryToUtf8 tests the Utf8() payload conversion method
// with a binary payload
func TestConvertBinaryToUtf8(t *testing.T) {
	payload := Payload{
		Encoding: Binary,
		Data:     []byte("ABC ёжз φπμλβωϘ"),
	}

	result, err := payload.Utf8()
	require.NoError(t, err)
	require.Equal(t, "ABC ёжз φπμλβωϘ", string(result))
	require.Len(t, result, 25)
}

// TestConvertUtf16ToUtf8 tests the Utf8() payload conversion method
// with a UTF16 encoded payload
func TestConvertUtf16ToUtf8(t *testing.T) {
	payload := Payload{
		Encoding: Utf16,
		Data: []byte{
			/* 0xFF 0xFE */ // byte order mark
			0x41, 0x00,
			0x42, 0x00,
			0x43, 0x00,
			0x20, 0x00,
			0x51, 0x04,
			0x36, 0x04,
			0x37, 0x04,
			0x20, 0x00,
			0xC6, 0x03,
			0xC0, 0x03,
			0xBC, 0x03,
			0xBB, 0x03,
			0xB2, 0x03,
			0xC9, 0x03,
			0xD8, 0x03,
		},
	}

	result, err := payload.Utf8()
	require.NoError(t, err)
	require.Equal(t, "ABC ёжз φπμλβωϘ", string(result))
	require.Len(t, result, 25)
}

// TestConvertCorruptUtf16 tests the Utf8() payload conversion method
// with a corrupted UTF16 payload
func TestConvertCorruptUtf16(t *testing.T) {
	payload := Payload{
		Encoding: Utf16,
		Data:     []byte{65, 66, 67}, // Odd number of bytes
	}

	result, err := payload.Utf8()
	require.Error(t, err)
	require.Len(t, result, 0)
}
