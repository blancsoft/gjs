package gjs

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

		defer recoverPanics(reject)
		jsFunc(resolve, reject)
		return nil
	})

	return promise.New(handler)
}

func reflectionOf(v js.Value, t reflect.Type) (reflect.Value, error) {
	// TODO: Implement this
	return reflect.Value{}, nil
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

func recoverPanics(callback any) {
	caller := GetCallerInfo(1)
	warn := console.Get("warn")
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
