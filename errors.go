package validator

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	err     error
	context string
	value   string
}

var ErrEmptyTag = NewError("empty tag")
var ErrEmptyTagValue = NewError("empty tag value")
var ErrInvalidKeyword = NewError("invalid keyword")
var ErrInvalidInteger = NewError("invalid integer value")
var ErrMissingEnumValue = NewError("missing enum values")

func NewError(msg string) *ValidationError {
	return &ValidationError{
		err: errors.New(msg),
	}
}

func (e *ValidationError) Context(context string) *ValidationError {
	if e == nil {
		return nil
	}

	e2 := e.Clone()
	e2.context = context

	return e2
}

func (e *ValidationError) Value(value any) *ValidationError {
	if e == nil {
		return nil
	}

	e2 := e.Clone()
	e2.context = fmt.Sprintf("%v", value)

	return e2
}

func (e *ValidationError) Clone() *ValidationError {
	if e == nil {
		return nil
	}

	return &ValidationError{
		err:     e.err,
		context: e.context,
		value:   e.value,
	}
}

func (e *ValidationError) Error() string {
	if e == nil {
		return "Success"
	}

	result := e.err.Error()
	if e.context != "" {
		result += ": " + e.context
	}

	return result
}
