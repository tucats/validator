package validator

import "encoding/json"

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

	return &item, nil
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
