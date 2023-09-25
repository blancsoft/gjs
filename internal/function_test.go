//go:build js && wasm

package internal_test

import (
	"syscall/js"
	"testing"

	. "github.com/chumaumenze/gjs/internal"
)

func add(a, b int) int {
	return a + b
}

func addVariadic(a int, b ...int) int {
	total := a
	for _, v := range b {
		total += v
	}
	return total
}

func TestFuncOf(t *testing.T) {
	fn, release := FuncOf(add)
	defer js.ValueOf(release).Invoke()

	js.Global().Set("add", fn)
	rv := js.Global().Call("eval", "add(1, 2)")

	t.Run("basic function", func(t *testing.T) {
		resolved, rejected := Await(rv)
		DeepEqualTest(t, true, rejected == nil)
		DeepEqualTest(t, 1, len(resolved))
		DeepEqualTest(t, 3, resolved[0].Int())
	})
}

func TestFuncOfVariadic(t *testing.T) {
	fn, release := FuncOf(addVariadic)
	defer js.ValueOf(release).Invoke()

	js.Global().Set("addVariadic", fn)
	rv := js.Global().Call("eval", "addVariadic(1, 2, 3, 4, 5, 6)")

	t.Run("variadic function", func(t *testing.T) {
		resolved, rejected := Await(rv)

		DeepEqualTest(t, true, rejected == nil)
		DeepEqualTest(t, 1, len(resolved))
		DeepEqualTest(t, 21, resolved[0].Int())
	})
}
