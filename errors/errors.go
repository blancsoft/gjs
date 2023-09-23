package errors

import (
	"fmt"
	"reflect"
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
