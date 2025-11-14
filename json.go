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

	b, _ := json.MarshalIndent(i, "", "   ")

	return string(b)
}

// NewJSON converts a JSON byt array into a validator item. IF the JSON did
// not contain a valid validator item, it will return an error.
func NewJSON(data []byte) (*Item, error) {
	var item Item

	err := json.Unmarshal(data, &item)
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

// NewFromJSON converts a JSON string into a validator item. If the JSON
// did not contain a valid validator item, it will return an error.
// This is different from NewJSON only in that this function accepts
// the JSON as a string instead of a byte array.
func NewFromJSON(data string, value any) (*Item, error) {
	// Parse a JSON string and convert it to an Item structure.
	v := Item{}

	err := json.Unmarshal([]byte(data), &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
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
