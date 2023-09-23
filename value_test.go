//go:build js && wasm

package gjs_test

import (
	"fmt"
	"syscall/js"
	"testing"

	. "github.com/chumaumenze/gjs"
	"github.com/stretchr/testify/assert"
)

func TestConvertToJSValue(t *testing.T) { //nolint:funlen
	type TestStruct struct {
		Field1 string `json:"field1"`
		Field2 int
	}

	// Test cases
	testCases := []struct {
		Name     string
		Input    interface{}
		Expected js.Value
	}{
		// Basic types
		{
			Name:     "NilValue",
			Input:    nil,
			Expected: js.Null(),
		},
		{
			Name:     "BoolValue",
			Input:    true,
			Expected: js.ValueOf(true),
		},
		{
			Name:     "IntValue",
			Input:    42,
			Expected: js.ValueOf(42),
		},
		{
			Name:     "StringValue",
			Input:    "Hello, World!",
			Expected: js.ValueOf("Hello, World!"),
		},
		{
			Name:     "FloatValue",
			Input:    3.14,
			Expected: js.ValueOf(3.14),
		},
		{
			Name:     "SliceValue",
			Input:    []int{1, 2, 3},
			Expected: js.ValueOf([]interface{}{1, 2, 3}),
		},
		{
			Name:     "MapValue",
			Input:    map[string]interface{}{"key": "value"},
			Expected: js.ValueOf(map[string]interface{}{"key": "value"}),
		},
		{
			Name: "StructValue",
			Input: TestStruct{
				Field1: "value1",
				Field2: 42,
			},
			Expected: js.ValueOf(map[string]interface{}{
				"field1": "value1",
				"Field2": 42,
			}),
		},
		{
			Name:     "Complex64Value",
			Input:    complex(float32(1.5), float32(2.5)),
			Expected: js.ValueOf(map[string]interface{}{"real": float32(1.5), "imag": float32(2.5)}),
		},
		{
			Name:     "Complex128Value",
			Input:    complex(2.0, 3.0),
			Expected: js.ValueOf(map[string]interface{}{"real": 2.0, "imag": 3.0}),
		},
		{
			Name:     "NilPointer",
			Input:    (*int)(nil),
			Expected: js.Null(),
		},
		{
			Name:     "PointerValue",
			Input:    func() *int { v := 42; return &v }(),
			Expected: js.ValueOf(42),
		},
		{
			Name:     "NilInterface",
			Input:    (interface{})(nil),
			Expected: js.Null(),
		},
		{
			Name:     "InterfaceValue",
			Input:    (interface{})(42),
			Expected: js.ValueOf(42),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := ValueOf(tc.Input)

			isEqual := deepEqual(actual, tc.Expected)
			if !isEqual {
				stringify := js.Global().Get("JSON").Get("stringify")
				errMsg := fmt.Sprintf("Comparison mismatch:\n\tx:  %s(%s)\n\tx:  %s(%s)\n",
					actual.Type().String(), stringify.Invoke(actual).String(),
					tc.Expected.Type().String(), stringify.Invoke(tc.Expected).String(),
				)
				fmt.Println(errMsg) //nolint:forbidigo
			}
			assert.True(t, isEqual)
		})
	}
}

func TestPromisify(t *testing.T) {
	expectedMessage := "promise resolved"
	handlerCalled := false
	handler := func(resolve, reject js.Value) {
		handlerCalled = true
		resolve.Invoke(expectedMessage)
	}
	Promisify(handler)

	t.Run("promise handler was called", func(t *testing.T) {
		assert.True(t, handlerCalled)
	})
}

func deepEqual(x, y js.Value) bool {
	if x.Type() != y.Type() {
		return false
	}

	if x.Type() != js.TypeObject {
		return x.Equal(y)
	}
	return compareObject(x, y)
}

func compareArray(x, y js.Value) bool {
	arrayConstructor := js.Global().Get("Array")
	isArrX := x.InstanceOf(arrayConstructor)
	isArrY := y.InstanceOf(arrayConstructor)
	if !(isArrX && isArrY) {
		return false
	}
	if x.Length() != y.Length() {
		return false
	}
	for i := 0; i < x.Length(); i++ {
		if !deepEqual(x.Index(i), y.Index(i)) {
			return false
		}
	}
	return true
}

func compareObject(x, y js.Value) bool {
	if x.Type() != js.TypeObject && x.Type() != y.Type() {
		return false
	}

	if x.InstanceOf(js.Global().Get("Array")) {
		return compareArray(x, y)
	}

	getKeys := js.Global().Get("Object").Get("keys")
	xKeys := getKeys.Invoke(x)
	yKeys := getKeys.Invoke(y)

	if xKeys.Length() == yKeys.Length() {
		for i := 0; i < xKeys.Length(); i++ {
			key := xKeys.Index(i).String()
			if !deepEqual(x.Get(key), y.Get(key)) {
				return false
			}
		}

		return true
	}

	return false
}
