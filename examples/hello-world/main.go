//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/chumaumenze/gjs"
)

type Data struct {
	Code    int
	Message string `json:"message"`
	inner   any
}

func main() {
	data := Data{
		Code:    200,
		Message: "Hello World!",
		inner:   "I am ignored",
	}
	_ = gjs.ValueOf(data) // e.g. {Code: 200, "message": "Hello World!"}

	find := func(nums []int, targets ...int) (idx []int, err error) {
		// ... implements find
		return
	}

	js.Global().Set("find", gjs.ValueOf(find))
	indexes := js.Global().Call("eval", `(async () => {
		return find([1,2,3,4,5], 1, 2, 4)
	})()`)

}
