package internal

import (
	"reflect"
	"syscall/js"
)

func StructOf(v any) js.Value {
	switch v := v.(type) {
	case reflect.Value:
		return structOf(v)
	default:
		return structOf(reflect.ValueOf(v))
	}
}

func structOf(v reflect.Value) js.Value {
	if v.Kind() != reflect.Struct {
		panic(ValueError{Type: reflect.Struct, Got: v.Kind()})
	}
	jsObject := Object.New()
	for i := 0; i < v.NumField(); i++ {
		// Ignore unexported fields
		if field := v.Field(i); field.CanInterface() {
			sf := v.Type().Field(i)
			jsObject.Set(nameOf(sf), valueOf(v.Field(i)))
		}
	}
	return jsObject
}
