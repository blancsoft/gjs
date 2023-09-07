//go:build js && wasm

package gjs

import (
	"reflect"
	"strings"
	"syscall/js"
)

var (
	Global  = js.Global()
	Array   = Global.Get("Array")
	Object  = Global.Get("Object")
	Console = Global.Get("console")
	Promise = Global.Get("Promise")
)

func ValueOf(v any) Value {
	switch v.(type) {
	case nil, js.Value:
		return Value(js.ValueOf(v))
	default:
		v := reflect.ValueOf(v)
		return Value(valueOf(v))
	}
}

func valueOfSlice(v reflect.Value) js.Value {
	length := v.Len()
	jsArray := Array.New()

	for i := 0; i < length; i++ {
		v := valueOf(v.Index(i))
		jsArray.SetIndex(i, v)
	}

	return jsArray
}

func valueOfMap(v reflect.Value) js.Value {
	jsObject := Object.New()
	for _, key := range v.MapKeys() {
		jsObject.Set(key.String(), valueOf(v.MapIndex(key)))
	}
	return jsObject
}

func valueOfStruct(v reflect.Value) js.Value {
	jsObject := Object.New()
	for i := 0; i < v.NumField(); i++ {
		// Ignore unexported fields
		if field := v.Field(i); field.CanInterface() {
			sf := v.Type().Field(i)

			// Use JSON tag if specified, otherwise use field name
			name := sf.Name
			if jsonTag := sf.Tag.Get("json"); jsonTag != "" {
				name = strings.SplitN(jsonTag, ",", 2)[0]
			}

			jsObject.Set(name, valueOf(v.Field(i)))
		}
	}
	return jsObject
}

func valueOfComplex(v reflect.Value) js.Value {
	c := v.Complex()
	jsObject := Object.New()
	jsObject.Set("real", real(c))
	jsObject.Set("imag", imag(c))
	return jsObject
}

func valueOfPointer(v reflect.Value) js.Value {
	if v.IsNil() {
		return js.Null()
	}
	return valueOf(v.Elem())
}

func valueOf(v reflect.Value) js.Value {
	switch v.Kind() { //nolint:exhaustive
	case reflect.Invalid:
		return js.Undefined()
	case reflect.Slice, reflect.Array:
		return valueOfSlice(v)
	case reflect.Map:
		return valueOfMap(v)
	case reflect.Struct:
		return valueOfStruct(v)
	case reflect.Complex64, reflect.Complex128:
		return valueOfComplex(v)
	case reflect.Interface, reflect.Pointer:
		return valueOfPointer(v)
	// case reflect.Func:
	//	return valueOfFunc(v)
	//	panic(errors.New("not implemented"))
	default:
		return js.ValueOf(v.Interface())
	}
}

type PromiseHandlerFunc func(resolve, reject js.Value)

func Promisify(jsFunc PromiseHandlerFunc) js.Value {
	handler := js.FuncOf(func(this js.Value, args []js.Value) any {
		resolve := args[0]
		reject := args[1]
		jsFunc(resolve, reject)
		return nil
	})

	return Promise.New(handler)
}
