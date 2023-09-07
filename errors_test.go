package gjs_test

import (
	"syscall/js"
	"testing"

	. "github.com/chumaumenze/gjs"
	"github.com/stretchr/testify/assert"
)

func TestRecoverPanicsWithGoFunc(t *testing.T) {
	expectedMsg := "panicking!"
	callbackCalled := false

	defer RecoverPanics(func(errMsg string) {
		callbackCalled = true
		t.Run("panic handler was called", func(t *testing.T) {
			assert.True(t, callbackCalled)
		})

		t.Run("error message matches", func(t *testing.T) {
			assert.Equal(t, errMsg, expectedMsg)
		})
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
