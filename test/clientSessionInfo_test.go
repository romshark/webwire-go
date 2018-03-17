package test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestClientSessionInfo tests the client.SessionInfo getter method
func TestClientSessionInfo(t *testing.T) {
	expectedBool := true
	expectedString := "somesamplestring1234"
	expectedInt := uint32(404)
	expectedNumber := float64(7.62)
	expectedArray := []string{"first", "second"}
	expectedStruct := struct {
		SampleString string `json:"struct_string"`
	}{
		SampleString: "sample struct string value",
	}

	// Initialize webwire server
	_, addr := setupServer(
		t,
		webwire.ServerOptions{
			Hooks: webwire.Hooks{
				OnRequest: func(ctx context.Context) (webwire.Payload, error) {
					// Extract request message and requesting client from the context
					msg := ctx.Value(webwire.Msg).(webwire.Message)

					// Try to create a new session
					if err := msg.Client.CreateSession(struct {
						SampleBool   bool     `json:"bool"`
						SampleString string   `json:"string"`
						SampleInt    uint32   `json:"int"`
						SampleNumber float64  `json:"number"`
						SampleArray  []string `json:"array"`
						SampleStruct struct {
							SampleString string `json:"struct_string"`
						} `json:"struct"`
					}{
						SampleBool:   expectedBool,
						SampleString: expectedString,
						SampleInt:    expectedInt,
						SampleNumber: expectedNumber,
						SampleArray:  expectedArray,
						SampleStruct: struct {
							SampleString string `json:"struct_string"`
						}{
							SampleString: expectedStruct.SampleString,
						},
					}); err != nil {
						return webwire.Payload{}, webwire.ReqErr{
							Code:    "INTERNAL_ERROR",
							Message: fmt.Sprintf("Internal server error: %s", err),
						}
					}
					return webwire.Payload{}, nil
				},
				// Define dummy hooks to enable sessions on this server
				OnSessionCreated: func(_ *webwire.Client) error { return nil },
				OnSessionLookup:  func(_ string) (*webwire.Session, error) { return nil, nil },
				OnSessionClosed:  func(_ *webwire.Client) error { return nil },
			},
		},
	)

	// Initialize client
	client := webwireClient.NewClient(
		addr,
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
	)
	defer client.Close()

	if err := client.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send authentication request and await reply
	if _, err := client.Request(
		"login",
		webwire.Payload{Data: []byte("credentials")},
	); err != nil {
		t.Fatalf("Request failed: %s", err)
	}

	// Verify getting inexistent field
	inexistent := client.SessionInfo("inexistent")
	if inexistent != nil {
		t.Fatalf("Expected nil for inexistent session info field, got: %v", inexistent)
	}

	// Verify field: bool
	samplebool, ok := client.SessionInfo("bool").(bool)
	if !ok {
		t.Fatalf(
			"Expected field 'bool' to be a boolean, got: %v",
			reflect.TypeOf(client.SessionInfo("bool")),
		)
	}
	if samplebool != expectedBool {
		t.Fatalf("Expected bool %t for field %s", expectedBool, "bool")
	}

	// Verify field: string
	samplestring, ok := client.SessionInfo("string").(string)
	if !ok {
		t.Fatalf(
			"Expected field 'string' to be a string, got: %v",
			reflect.TypeOf(client.SessionInfo("string")),
		)
	}
	if samplestring != expectedString {
		t.Fatalf("Expected string %s for field %s", expectedString, "string")
	}

	// Verify field: int
	sampleint, ok := client.SessionInfo("int").(float64)
	if !ok {
		t.Fatalf(
			"Expected field 'int' to be a float64, got: %v",
			reflect.TypeOf(client.SessionInfo("int")),
		)
	}
	if uint32(sampleint) != expectedInt {
		t.Fatalf("Expected uint32 (from float64) %d for field %s", expectedInt, "int")
	}

	// Verify field: number
	samplenumber, ok := client.SessionInfo("number").(float64)
	if !ok {
		t.Fatalf(
			"Expected field 'number' to be a float64, got: %v",
			reflect.TypeOf(client.SessionInfo("number")),
		)
	}
	if samplenumber != expectedNumber {
		t.Fatalf("Expected float64 number %f for field %s", expectedNumber, "number")
	}

	// Verify field: array
	samplearray, ok := client.SessionInfo("array").([]interface{})
	if !ok {
		t.Fatalf(
			"Expected field 'array' to be an array of empty interfaces, got: %v",
			reflect.TypeOf(client.SessionInfo("array")),
		)
	}
	for index, value := range samplearray {
		valStr, ok := value.(string)
		if !ok || expectedArray[index] != valStr {
			t.Fatalf(
				"Expected array item at index %d to be string('%s'), got: %v",
				index,
				expectedArray[index],
				value,
			)
		}
	}

	// Verify field: struct
	samplestruct, ok := client.SessionInfo("struct").(map[string]interface{})
	if !ok {
		t.Fatalf(
			"Expected field 'struct' to be a map of string -> empty interface, got: %v",
			reflect.TypeOf(client.SessionInfo("struct")),
		)
	}
	samplestructString, ok := samplestruct["struct_string"].(string)
	if !ok || samplestructString != expectedStruct.SampleString {
		t.Fatalf(
			"Expected struct field 'struct_string' to be a string('%s'), got: %v",
			expectedStruct.SampleString,
			samplestruct["struct_string"],
		)
	}
}
