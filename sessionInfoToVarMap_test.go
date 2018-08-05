package webwire

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSessionInfoToVarMap tests the SessionInfoToVarMap function
// using the generic session info implementation
func TestSessionInfoToVarMap(t *testing.T) {
	check := func(varMap map[string]interface{}) {
		expectedStruct := struct {
			Name   string
			Weight float64
		}{
			Name:   "samplename",
			Weight: 20.5,
		}

		require.Len(t, varMap, 4)
		require.Equal(t, "value1", varMap["field1"])
		require.Equal(t, int(42), varMap["field2"])
		require.Equal(t, expectedStruct, varMap["field3"])
		require.IsType(t, []string{}, varMap["field4"])
		require.ElementsMatch(t, []string{"item1", "item2"}, varMap["field4"])
	}

	info := &GenericSessionInfo{
		data: map[string]interface{}{
			"field1": "value1",
			"field2": int(42),
			"field3": struct {
				Name   string
				Weight float64
			}{
				Name:   "samplename",
				Weight: 20.5,
			},
			"field4": []string{"item1", "item2"},
		},
	}

	varMap := SessionInfoToVarMap(SessionInfo(info))
	check(varMap)

	// Test immutability, ensure fields won't mutate
	// even if the original session info object was changed
	info.data["field1"] = "mutated"
	info.data["field2"] = int(84)
	info.data["field3"] = struct {
		Name   string
		Weight float64
	}{
		Name:   "another name",
		Weight: 0.75,
	}
	info.data["field4"] = []string{"item3"}

	check(varMap)
}

// TestSessionInfoToVarMapNil tests the SessionInfoToVarMap function
// with a nil session info
func TestSessionInfoToVarMapNil(t *testing.T) {
	varMap := SessionInfoToVarMap(nil)
	require.Nil(t, varMap)
}
