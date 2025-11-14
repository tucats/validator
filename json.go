package validator

import (
	"encoding/json"
	"strings"
)

// String converts a validator item into a JSON-formatted string.
func (i *Item) String() string {
	if i == nil {
		return "<nil>"
	}

	// We do this is a two-step process. First, marshal the item as JSON.
	b, _ := json.MarshalIndent(i, "", "   ")

	// Now, read it back in into a map, where we can search for and manipulate
	// values that are numerically encoded to display as text.

	var m map[string]any

	err := json.Unmarshal(b, &m)
	if err != nil {
		return err.Error()
	}

	m = stringify(m, true)

	// Finally, convert the revised map to JSON and return it as a string.
	b, _ = json.MarshalIndent(m, "", "   ")

	return string(b)
}

func stringify(m map[string]any, toString bool) map[string]any {
	result := map[string]any{}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	for _, key := range keys {
		value := m[key]

		switch actual := value.(type) {
		case map[string]any:
			result[key] = stringify(actual, toString)

		case float64:
			if key == typeKeyName && toString {
				typeValue := Type(actual)

				result[key] = typeValue.String()
			} else {
				result[key] = actual
			}

		case int:
			if key == typeKeyName && toString {
				t := Type(actual)

				result[key] = t.String()
			} else {
				result[key] = actual
			}

		case string:
			if key == typeKeyName && !toString {
				if t, err := TypeFromString(actual); err == nil {
					result[key] = t
				}
			} else {
				result[key] = actual
			}

		case []any:
			newArray := []any{}

			for _, element := range actual {
				if m, ok := element.(map[string]any); ok {
					newArray = append(newArray, stringify(m, toString))
				} else {
					newArray = append(newArray, element)
				}
			}

			result[key] = newArray

		default:
			result[key] = value
		}
	}

	return result
}

// NewJSON converts a JSON byte array into a validator item. IF the JSON did
// not contain a valid validator item, it will return an error.
func NewJSON(data []byte) (*Item, error) {
	var (
		err  error
		item Item
		m    map[string]any
	)

	// First, unmarshal the JSON into a map.
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	// Search the map for "type" key values that are strings, that need to be converted
	// to int values.
	m = stringify(m, false)

	// Verify there are no misspelled or invalid keys
	if err := checkFields(m); err != nil {
		return nil, err
	}

	// Convert the map back to JSON
	data, err = json.Marshal(m)
	if err != nil {
		return nil, err
	}

	// Finally, unmarshal the JSON back into the validator item.
	err = json.Unmarshal(data, &item)
	if err != nil {
		return nil, err
	}

	return &item, check(&item)
}

// UnMarshal is a helper function that combines JSON validation and conversion of
// the representation into the destination structure.
//
// It will return errors if the validation tags on the structure definition are
// not correctly defines, if the JSON is not correctly formatted, or if the JSON
// contains fields or values that violate the validation rules.
//
// This is not recommended for use when the same structure is going to be validate
// many times, as it does not cache the validator structure.
func UnMarshal(data []byte, value any) error {
	// Validate the JSON against the specified data structure
	v, err := New(value)
	if err != nil {
		err = v.Validate(string(data))
		if err != nil {
			return json.Unmarshal(data, &value)
		}
	}

	return err
}

func checkFields(m map[string]any) error {
	// List of acceptable JSON field names for a validator item. IF
	// the json tag for the Item type is modified, it should be updated
	// in this list as well.
	fieldNames := map[string]bool{
		"name":              true,
		"alias":             true,
		"type":              true,
		"fields":            true,
		"min_value":         true,
		"has_min_value":     true,
		"max_value":         true,
		"has_max_value":     true,
		"base_type":         true,
		"min_length":        true,
		"has_min_length":    true,
		"max_length":        true,
		"has_max_length":    true,
		"enums":             true,
		"required":          true,
		"allow_foreign_key": true,
		"case_sensitive":    true,
	}

	// Verify all field names are valid
	for name, value := range m {
		if _, ok := fieldNames[name]; !ok {
			return ErrInvalidValidator.Context(name).Value("invalid field name")
		}

		if subMap, ok := value.(map[string]any); ok {
			if err := checkFields(subMap); err != nil {
				return err
			}
		}

		if subArray, ok := value.([]any); ok {
			for _, subValue := range subArray {
				if subMap, ok := subValue.(map[string]any); ok {
					if err := checkFields(subMap); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// check is an internal validator for validator structures.
func check(i *Item) error {
	if i == nil {
		return nil
	}

	// Check local items first.
	if i.ItemType == TypeInvalid || strings.HasPrefix(i.ItemType.String(), "unknown") {
		return ErrInvalidValidator.Context("type").Value("missing or invalid type")
	}

	// Check subordinate items.
	if err := check(i.BaseType); err != nil {
		return err
	}

	for _, field := range i.Fields {
		if err := check(field); err != nil {
			return err
		}
	}

	// If the min or max lengths are non-zero, verify that the flag is
	// set indicating they exit.
	if i.MinLength > 0 && !i.HasMinLength {
		return ErrInvalidValidator.Context("HasMinLength").Value("non-zero minLength without hasMinLength")
	}

	if i.MaxLength > 0 && !i.HasMaxLength {
		return ErrInvalidValidator.Context("HasMaxLength").Value("non-zero maxLength without hasMaxLength")
	}

	return nil
}
