package internal

import (
	"fmt"
	"reflect"
	"syscall/js"
)

type BaseError struct {
	Message   string
	Recovered any
}

func (e *BaseError) Error() string {
	if e.Message != "" {
		return e.Message
	} else {
		return fmt.Sprintf("unexpected panic: %+v", e.Recovered)
	}
}

type ValueError struct {
	BaseError
	Type reflect.Kind
	Got  reflect.Kind
}

func (e *ValueError) Error() string {
	if e.Message != "" {
		e.Message = fmt.Sprintf("expected %T, but go %T", e.Type, e.Got)
	}
	return e.Message
}

type ArgumentError struct {
	BaseError
	Expected int
	Got      int
	IsLess   bool
}

func (e *ArgumentError) Error() string {
	if e.Message != "" {
		s := "most"
		if e.IsLess {
			s = "least"
		}
		e.Message = fmt.Sprintf("expected at %s %d argument, got %d", s, e.Expected, e.Got)
	}
	return e.Message
}

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
