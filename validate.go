package validator

import (
	"encoding/json"
	"reflect"
	"strings"
)

// ValidateByName validates a JSON string against a named validator. If the
// named validator is not found, it returns an error. If the JSON string is
// valid according to the named validator, it returns nil.
func ValidateByName(name string, text string) error {
	item, exists := find(name)

	if !exists {
		return ErrUndefinedStructure.Context(name)
	}

	return item.Validate(text)
}

// For a given validator, determine if the provided JSON string is valid
// according to the rules for the given validator. If the JSON is not
// correctly formed or if a validation rule is violated, an error is
// returned. IF the provided JSON includes recursive or nested values
// that exceed the maximum recursion depth, an error is returned. If
// the JSON string is valid, it returns nil.
func (i *Item) Validate(text string) error {
	var (
		err error
		v   any
	)

	// Parse the JSON into an abstract object.
	err = json.Unmarshal([]byte(text), &v)
	if err != nil {
		return err
	}

	return i.validateValue(v, 0)
}

// This is the recursive validator function for a single item.
func (i *Item) validateValue(v any, depth int) error {
	if i == nil {
		return ErrNilValidator
	}

	if depth > maxValidationDepth {
		return ErrMaxDepthExceeded.Value(depth)
	}

	// If this item is an alias to another structure, resolve it now by
	// reading from the dictionary (using the reserved alias prefix).
	if i.Alias != "" {
		aliasItem, exists := find(aliasPrefix + i.Alias)
		if exists && aliasItem.ItemType == TypeStruct {
			i = aliasItem.Copy()
		}
	}

	// Based on the item's type, perform the appropriate validations.
	switch i.ItemType {
	case TypeAny:
		return nil // Accept anything.

	case TypePointer:
		return i.BaseType.validateValue(v, depth+1)

	case TypeArray:
		array, ok := v.([]any)
		if !ok {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

		if i.HasMinLength && len(array) < i.MinLength {
			return ErrArrayLengthOutOfRange.Context(i.Name).Value(len(array)).Expected(i.MinLength)
		}

		if i.HasMaxLength && len(array) > i.MaxLength {
			return ErrArrayLengthOutOfRange.Context(i.Name).Value(len(array)).Expected(i.MaxLength)
		}

		for _, element := range array {
			err := i.BaseType.validateValue(element, depth+1)
			if err != nil {
				return err
			}
		}

		return nil

	case TypeMap:
		// A current limitation of maps is that the key must always be of type string.
		if len(i.Enums) > 0 {
			actual := reflect.ValueOf(v)
			keys := actual.MapKeys()

			for _, key := range keys {
				// Validate that the key value itself is valid.
				found := false
				keyString := key.String()

				for _, enum := range i.Enums {
					if i.CaseSensitive {
						if keyString == enum {
							found = true

							break
						}
					} else {
						if strings.EqualFold(keyString, enum) {
							found = true

							break
						}
					}
				}

				if !found {
					return ErrInvalidEnumeratedValue.Context(i.Name).Value(keyString).Expected(i.Enums)
				}

				// Validate that the value of the key in this map is also valid.
				mapValue := actual.MapIndex(key).Interface()

				err := i.BaseType.validateValue(mapValue, depth+1)
				if err != nil {
					return err
				}
			}
		}

	case TypeTime:
		_, err := getTimeValue(v)
		if err != nil {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

		if i.HasMinValue {
			t, _ := getTimeValue(i.MinValue)
			vv, _ := getTimeValue(v)

			if vv.Before(t) {
				return ErrValueOutOfRange.Context(i.Name).Value(v)
			}
		}

		if i.HasMaxValue {
			t, _ := getTimeValue(i.MaxValue)
			vv, _ := getTimeValue(v)

			if vv.After(t) {
				return ErrValueOutOfRange.Context(i.Name).Value(v)
			}
		}

	case TypeDuration:
		vv, err := getDurationValue(v)
		if err != nil {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

		if i.HasMinValue {
			t, _ := getDurationValue(i.MinValue)

			if vv.Milliseconds() < t.Milliseconds() {
				return ErrValueOutOfRange.Context(i.Name).Value(v)
			}
		}

		if i.HasMaxValue {
			t, _ := getDurationValue(i.MaxValue)

			if vv.Milliseconds() > t.Milliseconds() {
				return ErrValueOutOfRange.Context(i.Name).Value(v)
			}
		}

	case TypeUUID:
		_, err := getUUIDValue(v)
		if err != nil {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

	case TypeBool:
		_, err := getBoolValue(v)
		if err != nil {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

	case TypeInt:
		value, err := getIntValue(v)
		if err != nil {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

		if i.HasMinValue {
			t, _ := getIntValue(i.MinValue)
			if value < t {
				return ErrValueOutOfRange.Context(i.Name).Value(value)
			}
		}

		if i.HasMaxValue {
			t, _ := getIntValue(i.MaxValue)
			if value > t {
				return ErrValueOutOfRange.Context(i.Name).Value(value)
			}
		}

		found := false

		for _, enum := range i.Enums {
			enumValue, _ := getIntValue(enum)

			if value == enumValue {
				found = true

				break
			}
		}

		if len(i.Enums) > 0 && !found {
			return ErrInvalidEnumeratedValue.Context(i.Name).Value(value).Expected(i.Enums)
		}

	case TypeFloat:
		value, err := getFloatValue(v)
		if err != nil {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

		if i.HasMinValue {
			t, _ := getFloatValue(i.MinValue)
			if value < t {
				return ErrValueOutOfRange.Context(i.Name).Value(value)
			}
		}

		if i.HasMaxValue {
			t, _ := getFloatValue(i.MaxValue)
			if value > t {
				return ErrValueOutOfRange.Context(i.Name).Value(value)
			}
		}

	case TypeList:
		value, err := getStringValue(v)
		if err != nil {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

		// Convert the string into a slice of strings. The only permitted separator is a comma.
		elements := strings.Split(value, ",")

		// IF there is a minimum length, check the length of the list
		if i.HasMinLength {
			if len(elements) < i.MinLength {
				return ErrValueLengthOutOfRange.Context(i.Name).Value(value)
			}
		}

		// If there is a maximum length, check the length of the list
		if i.HasMaxLength {
			if len(elements) > i.MaxLength {
				return ErrValueLengthOutOfRange.Context(i.Name).Value(value)
			}
		}

		// If there is an enum list, check each element in the list against the
		// enumerated values list. This may be case-sensitive depending on the
		// casematch setting for the item.
		if len(i.Enums) > 0 {
			for _, element := range elements {
				found := false
				element = strings.TrimSpace(element)

				for _, enum := range i.Enums {
					if i.CaseSensitive {
						if element == enum {
							found = true

							break
						}
					} else {
						if strings.EqualFold(element, enum) {
							found = true

							break
						}
					}
				}

				if !found {
					return ErrInvalidEnumeratedValue.Context(i.Name).Value(element).Expected(i.Enums)
				}
			}
		}

	case TypeString:
		value, err := getStringValue(v)
		if err != nil {
			return ErrInvalidData.Context(i.Name).Value(v)
		}

		if i.HasMinLength {
			if len(value) < i.MinLength {
				return ErrValueLengthOutOfRange.Context(i.Name).Value(value)
			}
		}

		if i.HasMaxLength {
			if len(value) > i.MaxLength {
				return ErrValueLengthOutOfRange.Context(i.Name).Value(value)
			}
		}

		found := false

		for _, enum := range i.Enums {
			if i.CaseSensitive {
				if value == enum {
					found = true

					break
				}
			} else {
				if strings.EqualFold(value, enum) {
					found = true

					break
				}
			}
		}

		if len(i.Enums) > 0 && !found {
			return ErrInvalidEnumeratedValue.Context(i.Name).Value(value).Expected(i.Enums)
		}

	case TypeStruct:
		m, ok := v.(map[string]any)
		if !ok {
			a, ok := v.([]map[string]any)
			if !ok || len(a) == 0 {
				return ErrInvalidData.Context(i.Name).Value(v)
			}

			m = a[0]
		}

		// Check if there are any field names that are not defined for the struct.
		// We do not do this if the validation allows "foreign" key values
		if !i.AllowForeignKey {
			for key := range m {
				found := false

				for _, field := range i.Fields {
					if field.Name == key {
						found = true

						break
					}
				}

				if !found {
					return ErrInvalidFieldName.Context(i.Name).Value(key)
				}
			}
		}

		// Verify each field found in the map against the struct's fields.
		for _, field := range i.Fields {
			fieldValue, exists := m[field.Name]
			if !exists {
				if field.Required {
					return ErrRequired.Value(field.Name)
				}

				continue
			}

			err := field.validateValue(fieldValue, depth+1)
			if err != nil {
				return err
			}
		}

	default:
		return ErrUnimplemented.Context(i.Name).Value(i.ItemType.String())
	}

	return nil
}
