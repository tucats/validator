package validator

import (
	"reflect"
	"strings"
)

const (
	// Maximum allowed recursion depth for validation. This is used to prevent infinite recursion.
	maxValidationDepth = 10

	// Prefix added to key values placed in the dictionary to represent a reference to
	// another validator. The user cannot create dictionary entries with this prefix.
	aliasPrefix = "_TYPE_ALIAS_"
)

// New accepts any value and returns a new validator for it. If the value is a structure,
// the validator will contain any rules specified in structure tags within the structure
// definition. Additional validation rules can be defined by calling Parse() with a tag
// string on the validator.
func New(v any) (*Item, error) {
	return defineItem(v, 0)
}

// defineItem handles defining a new validator for the given value. It can call itself
// recursively to a maximum allowed depth to handle nested structures, arrays of
// structures, etc.
func defineItem(v any, depth int) (*Item, error) {
	var err error

	// If we exceed maximum recursion depth, return an error
	if depth > maxValidationDepth {
		return nil, ErrMaxDepthExceeded.Value(depth)
	}

	// Based on the reflected type of this item, form an Item{} structure that defines it.
	item := &Item{}

	valueType := reflect.TypeOf(v)
	if valueType == nil {
		item.ItemType = TypeAny

		return item, nil
	}

	kind := valueType.Kind()

	// Handle well-known external types first. We have accessors, formatters, and
	// validators for these types even though they are external types, because they
	// are common value types found in JSON.
	typeName := valueType.String()
	switch typeName {
	case "uuid.UUID":
		item.ItemType = TypeUUID

		return item, nil

	case "time.Time":
		item.ItemType = TypeTime

		return item, nil

	case "time.Duration":
		item.ItemType = TypeDuration

		return item, nil
	}

	// Not one of the well-known package types, so handle based on the kind of the reflected type
	switch kind {
	case reflect.Interface:
		item.ItemType = TypeAny

	case reflect.Pointer:
		// Dereference the pointer and create an item for the base type
		valueType = valueType.Elem()
		v = reflect.Zero(valueType).Interface()
		item.ItemType = TypePointer

		item.BaseType, err = New(v)
		if err != nil {
			return nil, err
		}

	case reflect.Map:
		item.ItemType = TypeMap

	case reflect.String:
		item.ItemType = TypeString

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		item.ItemType = TypeInt

	case reflect.Float32, reflect.Float64:
		item.ItemType = TypeFloat

	case reflect.Bool:
		item.ItemType = TypeBool

	case reflect.Array, reflect.Slice:
		// Create an item for the base type of the array/slice
		baseItem, err := New(reflect.Zero(valueType.Elem()).Interface())
		if err != nil {
			return nil, err
		}

		item.ItemType = TypeArray
		item.BaseType = baseItem

	case reflect.Struct:
		var cacheThis bool

		// Iterate over the fields of the struct and build items for each field
		item.ItemType = TypeStruct

		// If the typename is a custom type, see if there is already an alias
		// for it. If so, reference the alias and we're done.
		if typeName != "struct" {
			previous, found := find(aliasPrefix + typeName)

			if found && previous.Alias != "" {
				item.Alias = typeName

				return item, nil
			}

			// Not already cached, let's create a shell in the dictionary for this
			// item now (to prevent infinite recursion) and set a flag to update it
			// when we've done this definition.
			cacheThis = true

			store(aliasPrefix+typeName, &Item{
				ItemType: TypeStruct,
				Alias:    typeName})
		}

		for i := 0; i < valueType.NumField(); i++ {
			field := valueType.Field(i)

			fieldItem, err := defineItem(reflect.Zero(field.Type).Interface(), depth+1)
			if err != nil {
				return nil, err
			}

			fieldItem.Name = field.Name

			// See if there is a JSON tag to get the field name from
			jsonTag := field.Tag.Get("json")
			if len(jsonTag) > 0 {
				jsonParts := strings.SplitN(jsonTag, ",", 2)
				if len(jsonParts) > 0 && jsonParts[0] != "" {
					fieldItem.Name = jsonParts[0]
				}
			}

			// Parse the field's validate tag if present and build an item for it
			tagString := field.Tag.Get(validateTagName)
			if len(strings.TrimSpace(tagString)) > 0 {
				err = fieldItem.ParseTag(tagString)
				if err != nil {
					return nil, err
				}
			}

			item.Fields = append(item.Fields, fieldItem)
		}

		if cacheThis {
			store(aliasPrefix+typeName, item)
		}

	default:
		err = ErrUnsupportedType.Value(valueType.String())
	}

	return item, err
}
