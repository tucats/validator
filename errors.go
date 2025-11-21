package validator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ValidationError represents a validation error. This includes the original error,
// the context of the validation (e.g., the field name), the actual value, and the
// expected values. If any of context, value, or expected values are empty, they
// are not included in the formatted error message string.
type ValidationError struct {
	err      error
	context  string
	value    string
	expected string
}

// Predefined validation errors.
var ErrArrayLengthOutOfRange = NewError("array length out of range")
var ErrEmptyTag = NewError("empty tag")
var ErrEmptyTagValue = NewError("empty tag value")
var ErrInvalidBaseTag = NewError("invalid base tag (only allowed on arrays and maps)")
var ErrInvalidData = NewError("invalid data")
var ErrInvalidDuration = NewError("invalid duration value")
var ErrInvalidEnumeratedValue = NewError("invalid enumerated value")
var ErrInvalidEnumType = NewError("invalid field type for enum, must be string or int")
var ErrInvalidFieldName = NewError("invalid field name")
var ErrInvalidInteger = NewError("invalid integer value")
var ErrInvalidKeyword = NewError("invalid keyword")
var ErrInvalidListTag = NewError("invalid list tag for item type")
var ErrInvalidName = NewError("invalid name")
var ErrInvalidTagName = NewError("invalid tag name")
var ErrInvalidValidator = NewError("invalid JSON instance of validator")
var ErrMaxDepthExceeded = NewError("maximum validation depth exceeded")
var ErrMissingEnumValue = NewError("missing enum values")
var ErrNameAlreadyExists = NewError("name already exists")
var ErrNilValidator = NewError("nil validator")
var ErrNotAMap = NewError("keyword only valid with map type")
var ErrRequired = NewError("required field missing")
var ErrSyntaxError = NewError("syntax error")
var ErrUndefinedStructure = NewError("undefined structure")
var ErrUnimplemented = NewError("unimplemented type")
var ErrUnsupportedType = NewError("unsupported type")
var ErrValueOutOfRange = NewError("value out of range")
var ErrValueLengthOutOfRange = NewError("value length out of range")

// Create a new validation error with the given message.
func NewError(msg string) *ValidationError {
	return &ValidationError{
		err: errors.New(msg),
	}
}

// Context adds a context value to the validation error. This
// returns the validation error with the updated context as
// a new validation error.
func (e *ValidationError) Context(context string) *ValidationError {
	if e == nil {
		return nil
	}

	e2 := e.copy()
	e2.context = context

	return e2
}

// Value adds a value to the validation error. This
// returns the validation error with the updated value as
// a new validation error. If the value is nil, the value
// field is not included in the formatted error message.
func (e *ValidationError) Value(value any) *ValidationError {
	if e == nil {
		return nil
	}

	if value == nil {
		e.value = ""

		return e
	}

	e2 := e.copy()
	e2.value = fmt.Sprintf("%v", value)

	return e2
}

// Expected adds expected values to the validation error. This is used
// when an enumeration rule is violated. The expected values can be a
// list of values, a slice of values, or a single value. This returns a
// copy of the validation error with the updated expected values as
// a new validation error.
func (e *ValidationError) Expected(expected ...any) *ValidationError {
	if e == nil {
		return nil
	}

	e2 := e.copy()

	list := make([]string, 0, len(expected))

	for _, v := range expected {
		switch v := v.(type) {
		case []string:
			for _, s := range v {
				list = append(list, s)
			}

			continue

		case []int:
			for _, n := range v {
				list = append(list, fmt.Sprintf("%d", n))
			}

		default:
			list = append(list, fmt.Sprintf("%v", v))
		}
	}

	e2.expected = strings.Join(list, ", ")
	if len(list) > 1 {
		e2.expected = "one of " + e2.expected
	}

	return e2
}

// Make a copy of the validation error. This allows the caller to use
// a predefined error variable, and add unique context, values, etc. to
// the message without modifying the original error.
func (e *ValidationError) copy() *ValidationError {
	if e == nil {
		return nil
	}

	return &ValidationError{
		err:     e.err,
		context: e.context,
		value:   e.value,
	}
}

// Error returns the formatted error message string. Supporting this
// function makes validation errors match the error interface in Go.
// The error is formatted with the original error message, context, value,
// and expected values. If any of context, value, or expected values are
// empty, they are not included in the formatted error message string.
func (e *ValidationError) Error() string {
	if e == nil {
		return "Success"
	}

	result := e.err.Error()
	if e.context != "" {
		result += ", in " + e.context
	}

	if e.value != "" {
		result += fmt.Sprintf(": %s", strconv.Quote(e.value))
	}

	if e.expected != "" {
		result += fmt.Sprintf(", expected %s", e.expected)
	}

	return result
}
