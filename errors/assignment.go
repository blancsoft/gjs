//go:build js && wasm
// +build js,wasm

package errors

import (
	"fmt"
	"reflect"
	"syscall/js"
)

type AssignmentError struct {
	BaseError

	Type js.Type
	Kind reflect.Kind
}

func (e *AssignmentError) Error() string {
	if e.Recovered != nil {
		return fmt.Sprintf("unexpected panic: %+v", e.Recovered)
	}
	if e.Type == js.TypeUndefined {
		return fmt.Sprintf("invalid assignment to Go kind: %v must be a non-nil pointer", e.Kind)
	}
	return fmt.Sprintf("invalid assignment from JS type: %v to Go kind: %v", e.Type, e.Kind)
}
