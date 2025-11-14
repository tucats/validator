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

	if item == nil {
		return ErrNilValidator.Context(tag)
	}

	// Split the tag string into parts by commas
	parts := Split(tag, ",")
	if len(parts) == 0 {
		return ErrEmptyTag
	}

	// Scan over each part, and apply it to the item as appropriate
	for _, part := range parts {
		var (
			key   string
			value string
		)

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
		case "type":
			switch value {
			case "string":
				item.ItemType = TypeString
			case "integer":
				item.ItemType = TypeInt
			case "float":
				item.ItemType = TypeFloat
			case "bool":
				item.ItemType = TypeBool
			case "array":
				item.ItemType = TypeArray
				item.BaseType = NewType(TypeAny)
			case "struct":
				item.ItemType = TypeStruct
				item.BaseType = NewType(TypeAny)
			case "pointer":
				item.ItemType = TypePointer
				item.BaseType = NewType(TypeAny)
			case "map":
				item.ItemType = TypeMap
				item.BaseType = NewType(TypeAny)

			default:
				return ErrUnsupportedType.Context(key).Value(value)
			}

		case "base", "value":
			if item.BaseType == nil {
				return ErrInvalidBaseTag.Context(tag)
			}

			// If the sub-tag is wrapped in parentheses or single quotes, remove them.
			if strings.HasPrefix(value, "(") && strings.HasSuffix(value, ")") {
				value = value[1 : len(value)-1]
			} else if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
				value = value[1 : len(value)-1]
			}

			// Parse the base type's tag and apply it to the BaseType in the current validator item.
			// IF the base type is an array or pointer, use the thing it points to.
			base := item.BaseType
			if base.ItemType == TypePointer || base.ItemType == TypeArray {
				base = base.BaseType
			}

			err = base.ParseTag(value)

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

		case "key":
			if item.ItemType != TypeMap {
				return ErrNotAMap.Context("key").Value(tag)
			}

			fallthrough

		case "enum", "enums":
			// Enum can't be used on a bool, struct, or pointer to struct
			if item.ItemType == TypeBool ||
				item.ItemType == TypeStruct ||
				item.ItemType == TypePointer && item.BaseType != nil && item.BaseType.ItemType == TypeStruct {
				return ErrInvalidEnumType.Context(key).Value(item.ItemType.String())
			}

			// The values could be separated by "|" characters, or they could be a nested list.
			sep := "|"

			if strings.HasPrefix(value, "(") && strings.HasSuffix(value, ")") {
				value = value[1 : len(value)-1]
				sep = ","
			} else if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
				value = value[1 : len(value)-1]
				sep = ","
			}

			// Make the list into an array of values.
			enums := strings.Split(value, sep)
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

// Split a string into separate components, using a defined separator character.
// If the separator is not provided, the function defaults to a comma ",".  The
// split ignores separators enclosed within single quotes or parentheses.
func Split(input string, separator string) []string {
	parts := make([]string, 0)
	current := ""
	inQuotes := false
	inDoubleQuotes := false
	inParens := 0
	sep := rune(',')

	for _, char := range separator {
		sep = rune(char)

		break
	}

	for _, char := range input {
		if char == '\'' {
			inQuotes = !inQuotes
		}

		if char == '"' {
			inDoubleQuotes = !inDoubleQuotes
		}

		if char == '(' {
			inParens++
		}

		if char == ')' && inParens > 0 {
			inParens--
		}

		if char == sep && inParens == 0 && !inDoubleQuotes && !inQuotes {
			parts = append(parts, strings.TrimSpace(current))
			current = ""
		} else {
			current += string(char)
		}
	}

	if current != "" {
		parts = append(parts, strings.TrimSpace(current))
	}

	return parts
}
