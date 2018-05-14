package test

import (
	"context"
	"reflect"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
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

// Copy implements the webwire.SessionInfo interface.
// It deep-copies the object and returns it's exact clone
func (sinf *testClientSessionInfoSessionInfo) Copy() webwire.SessionInfo {
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

// Fields implements the webwire.SessionInfo interface.
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

// Copy implements the webwire.SessionInfo interface.
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
) webwire.SessionInfo {
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
				clt *webwire.Client,
				_ *webwire.Message,
			) (webwire.Payload, error) {
				// Try to create a new session
				if err := clt.CreateSession(&testClientSessionInfoSessionInfo{
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
				}); err != nil {
					return webwire.Payload{}, err
				}
				return webwire.Payload{}, nil
			},
		},
		webwire.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
			SessionInfoParser:     testClientSessionInfoSessionInfoParser,
		},
		callbackPoweredClientHooks{},
	)
	defer client.connection.Close()

	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send authentication request and await reply
	if _, err := client.connection.Request(
		"login",
		webwire.Payload{Data: []byte("credentials")},
	); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify getting inexistent field
	inexistent := client.connection.SessionInfo("inexistent")
	if inexistent != nil {
		t.Fatalf(
			"Expected nil for inexistent session info field, got: %v",
			inexistent,
		)
	}

	// Verify field: bool
	samplebool, ok := client.connection.SessionInfo("bool").(bool)
	if !ok {
		t.Fatalf(
			"Expected field 'bool' to be boolean, got: %v",
			reflect.TypeOf(client.connection.SessionInfo("bool")),
		)
	}
	if samplebool != expectedBool {
		t.Fatalf("Expected bool %t for field %s", expectedBool, "bool")
	}

	// Verify field: string
	samplestring, ok := client.connection.SessionInfo("string").(string)
	if !ok {
		t.Fatalf(
			"Expected field 'string' to be string, got: %v",
			reflect.TypeOf(client.connection.SessionInfo("string")),
		)
	}
	if samplestring != expectedString {
		t.Fatalf("Expected string %s for field %s", expectedString, "string")
	}

	// Verify field: int
	sampleint, ok := client.connection.SessionInfo("int").(int)
	if !ok {
		t.Fatalf(
			"Expected field 'int' to be int, got: %v",
			reflect.TypeOf(client.connection.SessionInfo("int")),
		)
	}
	if sampleint != expectedInt {
		t.Fatalf(
			"Expected int %d for field %s",
			expectedInt,
			"int",
		)
	}

	// Verify field: number
	samplenumber, ok := client.connection.SessionInfo("number").(float64)
	if !ok {
		t.Fatalf(
			"Expected field 'number' to be float64, got: %v",
			reflect.TypeOf(client.connection.SessionInfo("number")),
		)
	}
	if samplenumber != expectedNumber {
		t.Fatalf(
			"Expected float64 number %f for field %s",
			expectedNumber,
			"number",
		)
	}

	// Verify field: array
	samplearray, ok := client.connection.SessionInfo("array").([]string)
	if !ok {
		t.Fatalf(
			"Expected field 'array' to be array of empty interfaces, got: %v",
			reflect.TypeOf(client.connection.SessionInfo("array")),
		)
	}
	for index, value := range samplearray {
		if expectedArray[index] != value {
			t.Fatalf(
				"Expected array item at index %d to be string('%s'), got: %v",
				index,
				expectedArray[index],
				value,
			)
		}
	}

	// Verify field: struct
	samplestruct, ok := client.connection.SessionInfo(
		"struct",
	).(testClientSessionInfoStruct)
	if !ok {
		t.Fatalf(
			"Expected field 'struct' to be map of empty interfaces, got: %v",
			reflect.TypeOf(client.connection.SessionInfo("struct")),
		)
	}
	if samplestruct.SampleString != expectedStruct.SampleString {
		t.Fatalf(
			"Expected struct field 'struct_string' to be string('%s'), got: %v",
			expectedStruct.SampleString,
			samplestruct.SampleString,
		)
	}
}
