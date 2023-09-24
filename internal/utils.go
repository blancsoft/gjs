package internal

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"syscall/js"
)

type PromiseHandlerFunc func(resolve, reject js.Value)

func Promisify(jsFunc PromiseHandlerFunc) js.Value {
	handler := js.FuncOf(func(this js.Value, args []js.Value) any {
		resolve := args[0]
		reject := args[1]

		defer RecoverPanics(reject)
		jsFunc(resolve, reject)
		return nil
	})

	return Promise.New(handler)
}

func reflectionOf(v js.Value, t reflect.Type) (out reflect.Value, err error) {
	out = reflect.New(t)
	err = recoverAssignTo(out, v)
	out = out.Elem()
	return
}

type CallerInfo struct {
	File       string
	LineNumber int
	FuncName   string
}

func (c *CallerInfo) String() string {
	i := strings.LastIndex(c.FuncName, ".")
	fn := c.FuncName[i+1 : len(c.FuncName)]

	return fmt.Sprintf("%s\t%s:%d", c.File, fn, c.LineNumber)
}

func RecoverPanics(callback any) {
	caller := GetCallerInfo(1)
	warn := Console.Get("warn")
	invalidArgMsg := caller.String() +
		"\tWARNING: Callback argument expects a Go/JS value function."

	if r := recover(); r != nil {
		errMsg := fmt.Sprintf("%+v", r)

		switch cb := callback.(type) {
		case js.Value:
			if cb.Type() != js.TypeFunction {
				warn.Invoke(invalidArgMsg)
				warn.Invoke(errMsg)
			} else {
				cb.Invoke(errMsg)
			}
		case func(string):
			cb(errMsg)
		default:
			println(invalidArgMsg) //nolint:forbidigo
			println(errMsg)        //nolint:forbidigo
		}
	}
}

func GetCallerInfo(skip int) CallerInfo {
	pc, file, lineNo, _ := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc).Name()

	return CallerInfo{
		File:       file,
		LineNumber: lineNo,
		FuncName:   funcName,
	}
}

func DeepEqual(x, y js.Value) bool {
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
		if !DeepEqual(x.Index(i), y.Index(i)) {
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
			if !DeepEqual(x.Get(key), y.Get(key)) {
				return false
			}
		}

		return true
	}

	return false
}

func Await(awaitable js.Value) ([]js.Value, []js.Value) {
	// Easy JS await implementation
	// Stolen from https://stackoverflow.com/a/68427221
	then := make(chan []js.Value)
	defer close(then)
	thenFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		then <- args
		return nil
	})
	defer thenFunc.Release()

	catch := make(chan []js.Value)
	defer close(catch)
	catchFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		catch <- args
		return nil
	})
	defer catchFunc.Release()

	awaitable.Call("then", thenFunc).Call("catch", catchFunc)

	select {
	case result := <-then:
		return result, nil
	case err := <-catch:
		return nil, err
	}
}

func DebugValue(v js.Value, depth int) string {
	insp := Global.Get("require").Invoke("util").Get("inspect")
	return insp.Invoke(v, false, depth, false).String()
}
