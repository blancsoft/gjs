package gjs

import (
	"reflect"
	"syscall/js"

	"github.com/chumaumenze/gjs/errors"
)

func FuncOf(v any) (jsFn js.Func, release js.Func) {
	switch v := v.(type) {
	case reflect.Value:
		return funcOf(v)
	default:
		return funcOf(reflect.ValueOf(v))
	}
}

func funcOf(v reflect.Value) (jsFn js.Func, release js.Func) {
	if v.Kind() != reflect.Func {
		panic(errors.ValueError{Type: reflect.Func, Got: v.Kind()})
	}

	vt := v.Type()
	jsFn = js.FuncOf(func(this js.Value, args []js.Value) any {
		return Promisify(func(resolve, reject js.Value) {
			err := errors.ArgumentError{Expected: vt.NumIn(), Got: len(args)}
			if len(args) < vt.NumIn() {
				err.IsLess = true
				reject.Invoke(err)
				return
			} else if !vt.IsVariadic() && len(args) > vt.NumIn() {
				reject.Invoke(err)
				return
			}

			var gargs []reflect.Value
			for _, a := range args {
				av, err := reflectionOf(a, v.Type())
				if err != nil {
					reject.Invoke(ValueOf(err))
					return
				}
				gargs = append(gargs, av)
			}
			if vt.IsVariadic() {
				posArgs := gargs[0 : len(gargs)-1]
				variadicArgs := gargs[len(gargs)-1].Interface().([]reflect.Value)
				gargs = append(posArgs, variadicArgs...)
			}

			go func() {
				rv := v.Call(gargs)
				resolve.Invoke(ValueOf(rv))
			}()
		})
	})

	release = js.FuncOf(func(this js.Value, args []js.Value) any {
		jsFn.Release()
		release.Release()
		return nil
	})

	return jsFn, release
}
