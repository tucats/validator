package validator

import (
	"reflect"
	"strings"
)

// Maximum allowed recursion depth for validation.
const (
	maxValidationDepth = 10
	aliasPrefix        = "_TYPE_ALIAS_"
)

func New(v any) (*Item, error) {
	return defineItem(v, 0)
}

func defineItem(v any, depth int) (*Item, error) {
	var err error

	// IF we exceed maximum recursion depth, return an error
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

	// Handle well-known external types first
	typeName := valueType.String()
	switch typeName {
	case "uuid.UUID":
		item.ItemType = TypeUUID

		return item, nil

	case "time.Time":
		item.ItemType = TypeTime

		return item, nil
	}

	// Handle based on the kind of the reflected type
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
