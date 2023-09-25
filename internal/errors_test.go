package internal_test

import (
	"syscall/js"
	"testing"

	. "github.com/chumaumenze/gjs/internal"
)

func TestRecoverPanicsWithGoFunc(t *testing.T) {
	var msg string
	expectedMsg := "panicking!"
	callbackCalled := false
	defer func() {
		DeepEqualTest(t, true, callbackCalled)
		DeepEqualTest(t, expectedMsg, msg)
	}()

	defer RecoverPanics(func(errMsg string) {
		callbackCalled = true
		msg = errMsg
	})
	panic(expectedMsg)
}

func TestRecoverPanicsWithJSValue(t *testing.T) {
	defer RecoverPanics(Console.Get("log"))
	panic("panicking WithJSValue")
}

func TestRecoverPanicsWithUnsupportedJSValue(t *testing.T) {
	defer RecoverPanics(js.Value{})
	panic("panicking WithUnsupportedJSValue")
}

func TestRecoverPanicsWithUnsupportedValue(t *testing.T) {
	defer RecoverPanics(nil)
	panic("panicking WithUnsupportedValue")
}
