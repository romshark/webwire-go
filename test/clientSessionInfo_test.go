package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	wwr "github.com/qbeon/webwire-go"
	wwrclt "github.com/qbeon/webwire-go/client"
)

type testClientSessionInfoStruct struct {
	SampleString string `json:"struct_string"`
}

type testClientSessionInfoSessionInfo struct {
	Bool   bool
	String string
	Int    int
	Number float64
	Array  []string
	Struct testClientSessionInfoStruct
}

// Copy implements the wwr.SessionInfo interface.
// It deep-copies the object and returns it's exact clone
func (sinf *testClientSessionInfoSessionInfo) Copy() wwr.SessionInfo {
	arrayClone := make([]string, len(sinf.Array))
	copy(arrayClone, sinf.Array)

	return &testClientSessionInfoSessionInfo{
		Bool:   sinf.Bool,
		String: sinf.String,
		Int:    sinf.Int,
		Number: sinf.Number,
		Array:  arrayClone,
		Struct: sinf.Struct,
	}
}

// Fields implements the wwr.SessionInfo interface.
// It returns a constant list of the names of all fields of the object
func (sinf *testClientSessionInfoSessionInfo) Fields() []string {
	return []string{
		"bool",
		"string",
		"int",
		"number",
		"array",
		"struct",
	}
}

// Copy implements the wwr.SessionInfo interface.
// It deep-copies the field identified by the provided name
// and returns it's exact clone
func (sinf *testClientSessionInfoSessionInfo) Value(
	fieldName string,
) interface{} {
	switch fieldName {
	case "bool":
		return sinf.Bool
	case "string":
		return sinf.String
	case "int":
		return sinf.Int
	case "number":
		return sinf.Number
	case "array":
		return sinf.Array
	case "struct":
		return sinf.Struct
	}
	return nil
}

func testClientSessionInfoSessionInfoParser(
	data map[string]interface{},
) wwr.SessionInfo {
	// Parse array field
	encodedArray := data["array"].([]interface{})
	typedArray := make([]string, len(encodedArray))
	for index := range encodedArray {
		typedArray[index] = encodedArray[index].(string)
	}

	// Parse struct field
	encodedStruct := data["struct"].(map[string]interface{})
	typedStruct := testClientSessionInfoStruct{
		SampleString: encodedStruct["struct_string"].(string),
	}

	return &testClientSessionInfoSessionInfo{
		Bool:   data["bool"].(bool),
		String: data["string"].(string),
		Int:    int(data["int"].(float64)),
		Number: data["number"].(float64),
		Array:  typedArray,
		Struct: typedStruct,
	}
}

// TestClientSessionInfo tests the client.SessionInfo getter method
func TestClientSessionInfo(t *testing.T) {
	expectedBool := true
	expectedString := "somesamplestring1234"
	expectedInt := int(404)
	expectedNumber := float64(7.62)
	expectedArray := []string{"first", "second"}
	expectedStruct := testClientSessionInfoStruct{
		SampleString: "sample struct string value",
	}

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				conn wwr.Connection,
				_ wwr.Message,
			) (wwr.Payload, error) {
				// Try to create a new session
				err := conn.CreateSession(&testClientSessionInfoSessionInfo{
					Bool:   expectedBool,
					String: expectedString,
					Int:    expectedInt,
					Number: expectedNumber,
					Array:  expectedArray,
					Struct: struct {
						SampleString string `json:"struct_string"`
					}{
						SampleString: expectedStruct.SampleString,
					},
				})
				assert.NoError(t, err)
				return nil, err
			},
		},
		wwr.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.AddressURL(),
		wwrclt.Options{
			DefaultRequestTimeout: 2 * time.Second,
			SessionInfoParser:     testClientSessionInfoSessionInfoParser,
		},
		callbackPoweredClientHooks{},
	)
	defer client.connection.Close()

	require.NoError(t, client.connection.Connect())

	// Send authentication request and await reply
	_, err := client.connection.Request(
		context.Background(),
		"login",
		wwr.NewPayload(wwr.EncodingBinary, []byte("credentials")),
	)
	require.NoError(t, err)

	// Verify getting inexistent field
	require.Nil(t, client.connection.SessionInfo("inexistent"))

	// Verify field: bool
	sampleBool := client.connection.SessionInfo("bool")
	require.IsType(t, expectedBool, sampleBool)
	require.Equal(t, expectedBool, sampleBool.(bool))

	// Verify field: string
	sampleString := client.connection.SessionInfo("string")
	require.IsType(t, expectedString, sampleString)
	require.Equal(t, expectedString, sampleString.(string))

	// Verify field: int
	sampleInt := client.connection.SessionInfo("int")
	require.IsType(t, expectedInt, sampleInt)
	require.Equal(t, expectedInt, sampleInt.(int))

	// Verify field: number
	sampleNumber := client.connection.SessionInfo("number")
	require.IsType(t, expectedNumber, sampleNumber)
	require.Equal(t, expectedNumber, sampleNumber.(float64))

	// Verify field: array
	sampleArray := client.connection.SessionInfo("array")
	require.IsType(t, expectedArray, sampleArray)
	require.Equal(t, expectedArray, sampleArray.([]string))

	// Verify field: struct
	sampleStruct := client.connection.SessionInfo("struct")
	require.IsType(t, expectedStruct, sampleStruct)
	require.Equal(t,
		expectedStruct.SampleString,
		sampleStruct.(testClientSessionInfoStruct).SampleString,
	)
}
