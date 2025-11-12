package validator

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/google/uuid"
)

// getTimeValue converts the given value to a time.Time if possible. If
// the provided value is not a time that can be converted to a time,
// return an error.
func getTimeValue(v any) (time.Time, error) {
	switch value := v.(type) {
	case time.Time:
		return value, nil

	case string:
		t, err := dateparse.ParseAny(value)
		if err != nil {
			return time.Time{}, ErrInvalidData.Value(value)
		}

		return t, nil

	default:
		return time.Time{}, ErrInvalidData.Value(value)
	}
}

// getDurationValue converts the given value to a time.Duration if possible.
// If the provided value is not a value that can be converted to a duration,
// return an error. This does allow an extension to the specification for a
// duration to allow days in the duration string, using a suffix of "d".
func getDurationValue(v any) (time.Duration, error) {
	switch value := v.(type) {
	case time.Duration:
		return value, nil

	case string:
		d, err := time.ParseDuration(value)
		if err != nil {
			return time.Duration(0), ErrInvalidData.Value(value)
		}

		return d, nil

	default:
		return time.Duration(0), ErrInvalidData.Value(value)
	}
}

// getUUIDValue converts the given value to a uuid.UUID if possible. If it
// is already a UUID, return it as is. If the value is a string or byte
// array, convert to a UUID. If the value is not a valid UUID, return an error.
func getUUIDValue(v any) (uuid.UUID, error) {
	switch value := v.(type) {
	case uuid.UUID:
		return value, nil

	case string:
		if value == "" {
			return uuid.Nil, nil
		}

		id, err := uuid.Parse(value)
		if err != nil {
			return uuid.Nil, ErrInvalidData.Value(value)
		}

		return id, nil

	case []byte:
		id, err := uuid.FromBytes(value)
		if err != nil {
			return uuid.Nil, ErrInvalidData.Value(value)
		}

		return id, nil

	default:
		return uuid.Nil, ErrInvalidData.Value(value)
	}
}

// getIntValue converts the given value to an int if possible. If it
// is already an int, return it as is. If the value is a string or float,
// convert it to an int. If the value is not a valid integer, return an error.
func getIntValue(v any) (int, error) {
	switch value := v.(type) {
	case int:
		return value, nil
	case float64:
		return int(value), nil
	case string:
		n, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return 0, ErrInvalidData.Value(value)
		}

		return int(n), nil

	default:
		return 0, ErrInvalidData.Value(value)
	}
}

// getFloatValue converts the given value to a float64 if possible. If it
// is already a float64, return it as is. If the value is a string or int,
// convert it to a float64. If the value is not a valid float, return an error.
func getFloatValue(v any) (float64, error) {
	switch value := v.(type) {
	case float64:
		return value, nil
	case int:
		return float64(value), nil

	case string:
		n, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, ErrInvalidData.Value(value)
		}

		return n, nil

	default:
		return 0, ErrInvalidData.Value(value)
	}
}

// getStringValue converts the given value to a string if possible. If it
// is a time, UUID, or duration, or float value, convert it to a string
// since these are all types that can be encoded as a string in JSON.
// OTherwise, return an error.
func getStringValue(v any) (string, error) {
	switch value := v.(type) {
	case string:
		return value, nil

	case float64:
		return fmt.Sprintf("%f", value), nil

	case time.Time:
		return value.Format(time.RFC3339), nil

	case time.Duration:
		return value.String(), nil

	case uuid.UUID:
		return value.String(), nil

	default:
		return "", ErrInvalidData.Value(value)
	}
}

// getBoolValue converts the given value to a bool if possible. If it
// is already a bool, return it as is. If the value is a string, convert
// it to a bool. If the value is not a valid boolean value, return an error.
func getBoolValue(v any) (bool, error) {
	switch value := v.(type) {
	case bool:
		return value, nil

	case string:
		switch strings.ToLower(value) {
		case "true":
			return true, nil
		case "false":
			return false, nil
		}

		return false, ErrInvalidData.Value(value)

	default:
		return false, ErrUnsupportedType
	}
}
