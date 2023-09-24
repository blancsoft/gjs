package internal

import (
	"reflect"
	"syscall/js"
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
		panic(ValueError{Type: reflect.Func, Got: v.Kind()})
	}

	vt := v.Type()
	jsFn = js.FuncOf(func(this js.Value, args []js.Value) any {
		return Promisify(func(resolve, reject js.Value) {
			err := ArgumentError{Expected: vt.NumIn(), Got: len(args)}
			if len(args) < vt.NumIn() {
				err.IsLess = true
				reject.Invoke(err.Error())
				return
			} else if !vt.IsVariadic() && len(args) > vt.NumIn() {
				reject.Invoke(err.Error())
				return
			}

			var variadicType reflect.Type
			if vt.IsVariadic() {
				variadicType = vt.In(vt.NumIn() - 1).Elem()
			}
			var gargs []reflect.Value
			for i, a := range args {
				var argType reflect.Type
				if vt.IsVariadic() && i >= vt.NumIn()-1 {
					argType = variadicType
				} else {
					argType = vt.In(i)
				}
				av, err := reflectionOf(a, argType)
				if err != nil {
					reject.Invoke(err.Error())
					return
				}
				gargs = append(gargs, av)
			}

			go func() {
				rv := v.Call(gargs)
				n := len(rv)
				if n > 0 {
					if isErrorType, err := hasError(rv[n-1]); isErrorType {
						rv = rv[:n-1]
						if err != nil {
							reject.Invoke(err.Error())
						}
					}
				}
				var jouts []any
				for _, x := range rv {
					jouts = append(jouts, ValueOf(x))
				}
				resolve.Invoke(jouts...)
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

func hasError(v reflect.Value) (bool, error) {
	if err, ok := v.Interface().(*error); ok {
		return true, *err
	}

	return false, nil
}
