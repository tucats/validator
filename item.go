package validator

import "fmt"

type Type int

const (
	TypeInvalid Type = iota
	TypeString
	TypeInt
	TypeFloat
	TypeBool
	TypeArray
	TypeStruct
	TypeUUID
	TypeTime
	TypeAny
)

const (
	validateTagName = "validate"
)

type Item struct {
	Name          string
	ValueType     Type
	ValueTypeName string
	BaseType      *Item
	MinLength     int
	MaxLength     int
	MinValue      any
	MaxValue      any
	Required      bool
	HasMinLength  bool
	HasMaxLength  bool
	HasMinValue   bool
	HasMaxValue   bool
	CaseSensitive bool
	IsPointer     bool
	IsArray       bool
	Enums         []string
	Fields        []Item
}

var Dictionary map[string]Item

func (t *Type) String() string {
	switch *t {
	case TypeString:
		return "string"
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeBool:
		return "bool"
	case TypeArray:
		return "array"
	case TypeStruct:
		return "struct"
	case TypeAny:
		return "any"
	case TypeUUID:
		return "uuid.UUID"
	default:
		return fmt.Sprintf("unknown type %d", t)
	}
}
