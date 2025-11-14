package validator

import (
	"fmt"
)

// Define an individual validator. This structure is used for any structure, field,
// map, array, or individual value. It contains flags representing all the validation
// rules defined for this item.
type Item struct {
	// The name of this item. This is usually the field name in a
	// structure. For other types, this is an empty string.
	Name string `json:"name,omitempty"`

	// IF this validator is a reference to another validator (or even a
	// references to the same validator), this will be the name of an
	// alias object in the dictionary. When set, the only other value
	// set in the validator is the ItemType value.
	Alias string `json:"alias,omitempty"`

	// The type of this item. This controls how the value in the JSON
	// payload being validated is interpreted, and controls what validations
	// are permitted for this item.
	ItemType Type `json:"type,omitempty"`

	// This is a list of the allowed values for this item. This is only
	// used for string, integer, and the key values for map types. IF there
	// are no enumerated values (enums), this will be an empty slice.
	Enums []string `json:"enums,omitempty"`

	// This is a list of the fields in the current structure. This is only
	// used for struct types. This will be an empty slice if there are no fields.
	Fields []*Item `json:"fields,omitempty"`

	// IF this is a pointer to another validator or an array of validators, this
	// is a pointer to the validator for the underlying type. If the current
	// validator is not a pointer or array, this will be nil.
	BaseType *Item `json:"base_type,omitempty"`

	// If there is a rule specifying a minimum length for a string, array, or
	// map key set, this will be the minimum length. The HasMinLength boolean
	// will be true if this value is used.
	MinLength int `json:"min_length,omitempty"`

	// If there is a rule specifying a maximum length for a string, array, or
	// map key set, this will be the maximum length. The HasMaxLength boolean
	// will be true if this value is used.
	MaxLength int `json:"max_length,omitempty"`

	// If there is a rule specifying a minimum value for a numeric field, this
	// will be the minimum value. The HasMinValue boolean will be true if
	// this value is used.
	MinValue any `json:"min_value,omitempty"`

	// If there is a rule specifying a maximum value for a numeric field, this
	// will be the maximum value. The HasMaxValue boolean will be true if
	// this value is used.
	MaxValue any `json:"max_value,omitempty"`

	// If there is a rule specifying that this field must always be present
	// in a json payload, this will be true.
	Required bool `json:"required,omitempty"`

	// If there is a rule specifying that the json can contain field names
	// not explicitly defined in the validator, this will be true. The default
	// is false, which means the json cannot contain field names not explicitly
	// defined in the validator.
	AllowForeignKey bool `json:"allow_foreign_key,omitempty"`

	// If there is a minimum length specified for this item, this is true.
	HasMinLength bool `json:"has_min_length,omitempty"`

	// If there is a maximum length specified for this item, this is true.
	HasMaxLength bool `json:"has_max_length,omitempty"`

	// If there is a minimum value specified for this item, this is true.
	HasMinValue bool `json:"has_min_value,omitempty"`

	// If there is a maximum value specified for this item, this is true.
	HasMaxValue bool `json:"has_max_value,omitempty"`

	// IF there is a rule specifying enumerated values, this is true when
	// the values are case-sensitive. By default, string values are not
	// case-sensitive.
	CaseSensitive bool `json:"case_sensitive,omitempty"`
}

const (
	typeKeyName = "type"
)

// SetRequired sets whether this item is required or not. By default, a json
// payload does not have to explicitly specify fields. If the field must be
// present in the JSON, then set this to true.
func (i *Item) SetRequired(required bool) *Item {
	if i == nil {
		return nil
	}

	i.Required = required

	return i
}

// SetMinValue sets the minimum allowed value for this item. The minimum
// value can be any numeric value.
func (i *Item) SetMinValue(v any) *Item {
	if i == nil {
		return nil
	}

	i.MinValue = v
	i.HasMinValue = true

	return i
}

// SetMaxValue sets the maximum allowed value for this item. The maximum
// value can be any numeric value.
func (i *Item) SetMaxValue(v any) *Item {
	if i == nil {
		return nil
	}

	i.MaxValue = v
	i.HasMaxValue = true

	return i
}

// SetMinLength sets the minimum allowed length for this item. The minimum
// length can be any integer value. This applies to string length, array
// length, or number of map keys.
func (i *Item) SetMinLength(l int) *Item {
	if i == nil {
		return nil
	}

	i.MinLength = l
	i.HasMinLength = true

	return i
}

// SetMaxLength sets the maximum allowed length for this item. The maximum
// length can be any integer value. This applies to string length, array
// length, or number of map keys.
func (i *Item) SetMaxLength(l int) *Item {
	if i == nil {
		return nil
	}

	i.MaxLength = l
	i.HasMaxLength = true

	return i
}

// SetEnums sets the allowed values for this item. The allowed values can be
// any string, integer, or any other type that can be converted to a string.
// This test can be applied to a string value, an array or strings or integers,
// or a string list where the values are comma-separated.
func (i *Item) SetEnums(values ...any) *Item {
	if i == nil {
		return nil
	}

	list := []string{}

	for _, v := range values {
		switch actual := v.(type) {
		case []string:
			list = append(list, actual...)

		case []int:
			for _, n := range actual {
				list = append(list, fmt.Sprintf("%d", n))
			}

		default:
			list = append(list, fmt.Sprintf("%v", actual))
		}
	}

	i.Enums = list

	return i
}

// SetName sets the field name for this item. The field name defaults to
// the field name in the struct tag. If the struct tag does not have a name
// field, it looks to see if there is a json tag and uses the name from that
// tag. Finally, it assumes the field name is the same as the variable name
// in the structure. This function can be used to override any of these defaults.
func (i *Item) SetName(name string) *Item {
	if i == nil {
		return nil
	}

	i.Name = name

	return i
}

// SetMatchCase sets the match case flag. By default, enumerated values are
// not case-sensitive when they are string values. If these are meant to be
// case-sensitive, set this flag to true.
func (i *Item) SetMatchCase(b bool) *Item {
	if i == nil {
		return nil
	}

	i.CaseSensitive = b

	return i
}

// SetForeignKeys sets the allow foreign key flag. By default, foreign keys
// (field names not defined in the validator) are not allowed. If the
// validator should instead ignore key values not defined in the validator,
// set this flag to true.
func (i *Item) SetForeignKeys(b bool) *Item {
	if i == nil {
		return nil
	}

	i.AllowForeignKey = b

	return i
}

// SetField adds a field validator to an existing structure validator. The field
// validator must have been previously completely defined (you cannot update a
// field after it is defined in the structure). The index is the zero-based index
// of the field in the structure.
//
// IF the validator is not for a structure or the index is negative, no change
// is made to the item.
func (i *Item) SetField(index int, v Item) *Item {
	if i == nil || i.ItemType != TypeStruct || index < 0 {
		return i
	}

	for len(i.Fields) <= index {
		i.Fields = append(i.Fields, &Item{})
	}

	i.Fields[index] = &v

	return i
}

// AddField adds a field validator to an existing structure validator. The field
// validator must have been previously completely defined (you cannot update a
// field after it is defined in the structure). The field is appended to the
// list of existing fields.
func (i *Item) AddField(v Item) *Item {
	if i == nil {
		return nil
	}

	i.Fields = append(i.Fields, &v)

	return i
}

// Copy creates a deep copy of the Item structure. The copied structure
// is completely independent of the original structure.
func (i *Item) Copy() *Item {
	if i == nil {
		return nil
	}

	result := &Item{
		Name:            i.Name,
		Alias:           i.Alias,
		ItemType:        i.ItemType,
		Enums:           append([]string{}, i.Enums...),
		Fields:          make([]*Item, len(i.Fields)),
		BaseType:        i.BaseType.Copy(),
		MinLength:       i.MinLength,
		MaxLength:       i.MaxLength,
		MinValue:        i.MinValue,
		MaxValue:        i.MaxValue,
		Required:        i.Required,
		AllowForeignKey: i.AllowForeignKey,
		HasMinLength:    i.HasMinLength,
		HasMaxLength:    i.HasMaxLength,
		HasMinValue:     i.HasMinValue,
		HasMaxValue:     i.HasMaxValue,
		CaseSensitive:   i.CaseSensitive,
	}

	for j, field := range i.Fields {
		result.Fields[j] = field.Copy()
	}

	return result
}
