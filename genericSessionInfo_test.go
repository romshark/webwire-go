package webwire

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGenericSessionInfoCopy tests the Copy method
// of the generic session info implementation
func TestGenericSessionInfoCopy(t *testing.T) {
	original := GenericSessionInfo{
		data: map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
		},
	}

	copied := original.Copy()

	check := func() {
		require.Equal(t, []string{"field1", "field2"}, copied.Fields())
		require.Equal(t, "value1", copied.Value("field1"))
		require.Equal(t, "value2", copied.Value("field2"))
	}

	// Verify consistency
	check()

	// Verify immutability
	delete(original.data, "field1")
	original.data["field2"] = "another_value"
	original.data["field3"] = "another_value"
	check()
}

// TestGenericSessionInfoValue tests the Value getter method
// of the generic session info implementation
func TestGenericSessionInfoValue(t *testing.T) {
	info := GenericSessionInfo{
		data: map[string]interface{}{
			"string":         "stringValue",
			"float":          12.5,
			"arrayOfStrings": []string{"item1", "item2"},
		},
	}

	// Check types
	require.IsType(t, float64(12.5), info.Value("float"))
	require.IsType(t, []string{}, info.Value("arrayOfStrings"))

	// Check values
	require.Equal(t, "stringValue", info.Value("string"))
	require.Equal(t, 12.5, info.Value("float"))
	require.Equal(t, []string{"item1", "item2"}, info.Value("arrayOfStrings"))

	// Check value of inexistent field
	require.Nil(t, info.Value("inexistent"))
}
