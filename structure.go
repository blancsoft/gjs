package gjs

import (
	"github.com/chumaumenze/gjs/errors"
	"reflect"
	"strings"
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
		panic(errors.ValueError{Type: reflect.Struct, Got: v.Kind()})
	}
	jsObject := object.New()
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
