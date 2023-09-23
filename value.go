//go:build js && wasm

package gjs

import (
	"reflect"
	"syscall/js"
)

var (
	global  = js.Global()
	array   = global.Get("Array")
	object  = global.Get("Object")
	console = global.Get("console")
	promise = global.Get("Promise")
	jsGo    = global.Get("Go")
)

func ValueOf(v any) js.Value {
	switch v := v.(type) {
	case nil, js.Value:
		return js.ValueOf(v)
	case reflect.Value:
		return valueOf(v)
	default:
		return valueOf(reflect.ValueOf(v))
	}
}

func valueOfFunc(v reflect.Value) js.Value {
	jsFn, release := funcOf(v)

	jsFnValue := js.ValueOf(jsFn)
	jsFnValue.Set("release", js.ValueOf(release))
	return jsFnValue
}

func valueOfSlice(v reflect.Value) js.Value {
	length := v.Len()
	jsArray := array.New()

	for i := 0; i < length; i++ {
		v := valueOf(v.Index(i))
		jsArray.SetIndex(i, v)
	}

	return jsArray
}

func valueOfMap(v reflect.Value) js.Value {
	jsObject := object.New()
	for _, key := range v.MapKeys() {
		jsObject.Set(key.String(), valueOf(v.MapIndex(key)))
	}
	return jsObject
}

func valueOfComplex(v reflect.Value) js.Value {
	c := v.Complex()
	jsObject := object.New()
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
		return structOf(v)
	case reflect.Complex64, reflect.Complex128:
		return valueOfComplex(v)
	case reflect.Interface, reflect.Pointer:
		return valueOfPointer(v)
	case reflect.Func:
		return valueOfFunc(v)
	default:
		return js.ValueOf(v.Interface())
	}
}
