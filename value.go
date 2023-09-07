//go:build js && wasm

package gjs

import "syscall/js"

type Value js.Value

func (v Value) Into() js.Value {
	return js.Value(v)
}
