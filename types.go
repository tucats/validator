package validator

import (
	"strconv"
)

// Define an int type that indicates what the underlying datatype of an item
// is. This is used to determine what kind of validation(s) can be done on a
// given item/field.

type Type int

// Enumerate all possible types that an item/field can have. Always add to the
// end to avoid breaking any existing JSON representation of a validator.
const (
	TypeInvalid Type = iota
	TypeString
	TypeInt
	TypeFloat
	TypeBool
	TypeStruct
	TypeArray
	TypePointer
	TypeMap // Any map type falls in this bucket. Only validations are for key values
	TypeUUID
	TypeTime
	TypeAny
	TypeList
	TypeDuration
)

// Map used to convert Type values to a string name.
var TypeNames = map[Type]string{
	TypeInvalid:  "invalid",
	TypeString:   "string",
	TypeInt:      "int",
	TypeFloat:    "float",
	TypeBool:     "bool",
	TypeStruct:   "struct",
	TypeAny:      "any",
	TypeUUID:     "uuid.UUID",
	TypeTime:     "time.Time",
	TypeDuration: "time.Duration",
	TypeMap:      "map[string]any",
	TypeList:     "stringList",
}

// String method for Type to return the string name of the type. Mostly
// used for debugging purposes.
func (t *Type) String() string {
	if name, ok := TypeNames[*t]; ok {
		return name
	}

	return "unknown type: " + strconv.Itoa(int(*t))
}

// NewType creates a new validator with the given type. There is no other
// information stored in the validator, so the caller should use ParseTag()
// or the explicit rule creation methods to add rules to the validator.
func NewType(kind Type) *Item {
	return &Item{ItemType: kind}
}
