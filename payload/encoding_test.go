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
