package validator

type Type int

const (
	TypeInvalid Type = iota
	TypeString
	TypeInt
	TypeFloat
	TypeBool
	TypeArray
	TypeStruct
	TypeFunction
	TypeAny
)

type Item struct {
	Name          string
	ValueType     Type
	ValueTypeName string
	Required      bool
	HasMinLength  bool
	MinLength     int
	HasMaxLength  bool
	MaxLength     int
	HasMinValue   bool
	MinValue      any
	HasMaxValue   bool
	MaxValue      any
	Enums         []string
	CaseSensitive bool
}

var Dictionary map[string]Item
