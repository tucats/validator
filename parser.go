package validator

import (
	"strconv"
	"strings"
)

// This is the structure tag name used to define validation rules for fields. The
// default can be overridden by the user before defining the first structure if needed.
var validateTagName = "validate"

// SetTagNAme sets the name of the tag used to define validation rules for fields.
// The default is "validate". If you need to change the tag string name, this must
// be called before using a function that reads a structure to create a validator.
func SetTagName(name string) error {
	if name == "" {
		return ErrInvalidTagName
	}

	validateTagName = name

	return nil
}

// ParseTag updates an existing validator to reflect the rules specified in
// the tag string. This tag string is normally extracted from struct field
// tags automatically when a new validator is created, but can also be used
// to set validator rules on a non-struct validator type.
func (item *Item) ParseTag(tag string) error {
	var err error

	// Split the tag string into parts by commas
	parts := strings.SplitSeq(tag, ",")

	// Scan over each part, and apply it to the item as appropriate
	for part := range parts {
		var (
			key   string
			value string
		)

		part = strings.TrimSpace(part)
		if part == "" {
			return ErrEmptyTag
		}

		// Elements can either be keywords (like "required") or key-value pairs
		// (like "minlength=5"). Break apart the item on the equals sign and make
		// sure there are the correct number of elements. You cannot specify a
		// key without a value. Key names are case-insensitive.
		elements := strings.SplitN(part, "=", 2)
		key = strings.ToLower(strings.TrimSpace(elements[0]))

		if len(elements) > 1 {
			value = strings.TrimSpace(elements[1])
			if len(value) == 0 {
				return ErrEmptyTagValue.Context(key)
			}
		}

		// Based on the key, apply the value to the item.
		switch key {
		case "list":
			if item.ItemType != TypeString {
				return ErrInvalidListTag.Context(key)
			}

			item.ItemType = TypeList

		case "name":
			item.Name = value

		case "required":
			item.Required = true

		case "minlength", "minlen":
			item.HasMinLength = true

			n, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return ErrInvalidInteger.Context(key).Value(value)
			}

			item.MinLength = int(n)

		case "maxlength", "maxlen":
			item.HasMaxLength = true

			n, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return ErrInvalidInteger.Context(key).Value(value)
			}

			item.MaxLength = int(n)

		case "maxvalue", "max":
			item.HasMaxValue = true

			item.MaxValue = value

		case "minvalue", "min":
			item.HasMinValue = true

			item.MinValue = value

		case "enum":
			// enum can only be used for strings, integers, maps, or lists
			if item.ItemType != TypeString && item.ItemType != TypeInt && item.ItemType != TypeMap && item.ItemType != TypeList {
				return ErrInvalidEnumType.Context(key).Value(item.ItemType.String())
			}

			// Value separator for the enumerated values is the pipe character "|".
			enums := strings.Split(value, "|")
			if len(enums) == 0 {
				return ErrMissingEnumValue.Context(key)
			}

			item.Enums = make([]string, len(enums))

			for i, enum := range enums {
				item.Enums[i] = strings.TrimSpace(enum)
			}

		case "matchcase", "casesensitive":
			item.CaseSensitive = true

		default:
			return ErrInvalidKeyword.Value(key)
		}
	}

	return err
}
