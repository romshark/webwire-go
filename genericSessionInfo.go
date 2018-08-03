package webwire

import (
	"reflect"
)

func deepCopy(src interface{}) interface{} {
	if src == nil {
		return nil
	}

	// Make the interface a reflect.Value
	original := reflect.ValueOf(src)

	// Make a copy of the same type as the original.
	cpy := reflect.New(original.Type()).Elem()

	// Recursively copy the original.
	copyRecursive(original, cpy)

	// Return the copy as an interface.
	return cpy.Interface()
}

func copyRecursive(original, cpy reflect.Value) {
	// handle according to original's Kind
	switch original.Kind() {
	case reflect.Interface:
		// If this is a nil, don't do anything
		if original.IsNil() {
			return
		}
		// Get the value for the interface, not the pointer.
		originalValue := original.Elem()

		// Get the value by calling Elem().
		copyValue := reflect.New(originalValue.Type()).Elem()
		copyRecursive(originalValue, copyValue)
		cpy.Set(copyValue)

	case reflect.Slice:
		if original.IsNil() {
			return
		}
		// Make a new slice and copy each element.
		cpy.Set(reflect.MakeSlice(
			original.Type(),
			original.Len(),
			original.Cap(),
		))
		for i := 0; i < original.Len(); i++ {
			copyRecursive(original.Index(i), cpy.Index(i))
		}

	case reflect.Map:
		if original.IsNil() {
			return
		}
		cpy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			copyRecursive(originalValue, copyValue)
			copyKey := deepCopy(key.Interface())
			cpy.SetMapIndex(reflect.ValueOf(copyKey), copyValue)
		}

	default:
		cpy.Set(original)
	}
}

// GenericSessionInfo defines a default webwire.SessionInfo interface
// implementation type used by the client when no explicit session info
// parser is used
type GenericSessionInfo struct {
	data map[string]interface{}
}

// Copy implements the webwire.SessionInfo interface.
// It deep-copies the object and returns it's exact clone
func (sinf *GenericSessionInfo) Copy() SessionInfo {
	return &GenericSessionInfo{
		data: deepCopy(sinf.data).(map[string]interface{}),
	}
}

// Fields implements the webwire.SessionInfo interface.
// It returns a constant list of the names of all fields of the object
func (sinf *GenericSessionInfo) Fields() []string {
	if sinf.data == nil {
		return make([]string, 0)
	}
	names := make([]string, len(sinf.data))
	index := 0
	for fieldName := range sinf.data {
		names[index] = fieldName
		index++
	}
	return names
}

// Value implements the webwire.SessionInfo interface.
// It returns an exact deep copy of a session info field value
func (sinf *GenericSessionInfo) Value(fieldName string) interface{} {
	if sinf.data == nil {
		return nil
	}
	if val, exists := sinf.data[fieldName]; exists {
		return deepCopy(val)
	}
	return nil
}

// GenericSessionInfoParser represents a default implementation of a
// session info object parser. It parses the info object into a generic
// session info type implementing the webwire.SessionInfo interface
func GenericSessionInfoParser(data map[string]interface{}) SessionInfo {
	return &GenericSessionInfo{data}
}
