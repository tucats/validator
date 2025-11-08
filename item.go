package validator

import (
	"encoding/json"
	"fmt"
)

type Type int

const (
	TypeInvalid Type = iota
	TypeString
	TypeInt
	TypeFloat
	TypeBool
	TypeStruct
	TypeMap // Any map type falls in this bucket. Only validations are for key values
	TypeUUID
	TypeTime
	TypeAny
)

const (
	validateTagName = "validate"
)

type Item struct {
	Name            string   `json:"name,omitempty"`
	ValueType       Type     `json:"type,omitempty"`
	Enums           []string `json:"enums,omitempty"`
	Fields          []Item   `json:"fields,omitempty"`
	BaseType        *Item    `json:"base_type,omitempty"`
	MinLength       int      `json:"min_length,omitempty"`
	MaxLength       int      `json:"max_length,omitempty"`
	MinValue        any      `json:"min_value,omitempty"`
	MaxValue        any      `json:"max_value,omitempty"`
	Required        bool     `json:"required,omitempty"`
	AllowForeignKey bool     `json:"allow_foreign_key,omitempty"`
	HasMinLength    bool     `json:"has_min_length,omitempty"`
	HasMaxLength    bool     `json:"has_max_length,omitempty"`
	HasMinValue     bool     `json:"has_min_value,omitempty"`
	HasMaxValue     bool     `json:"has_max_value,omitempty"`
	CaseSensitive   bool     `json:"case_sensitive,omitempty"`
	IsPointer       bool     `json:"is_pointer,omitempty"`
	IsArray         bool     `json:"is_array,omitempty"`
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
	case TypeStruct:
		return "struct"
	case TypeAny:
		return "any"
	case TypeUUID:
		return "uuid.UUID"
	case TypeTime:
		return "time.Time"
	case TypeMap:
		return "map[string]any"

	default:
		return fmt.Sprintf("unknown type %d", t)
	}
}

func NewType(kind Type) *Item {
	return &Item{ValueType: kind}
}

func (i *Item) SetMinValue(v any) *Item {
	i.MinValue = v
	i.HasMinValue = true

	return i
}

func (i *Item) SetMaxValue(v any) *Item {
	i.MaxValue = v
	i.HasMaxValue = true

	return i
}

func (i *Item) SetMinLength(l int) *Item {
	i.MinLength = l
	i.HasMinLength = true

	return i
}

func (i *Item) SetMaxLength(l int) *Item {
	i.MaxLength = l
	i.HasMaxLength = true

	return i
}

func (i *Item) SetEnums(values ...any) *Item {
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

func (i *Item) SetName(name string) *Item {
	i.Name = name

	return i
}

func (i *Item) SetMatchCase(b bool) *Item {
	i.CaseSensitive = b

	return i
}

func (i *Item) SetForeignKeys(b bool) *Item {
	i.AllowForeignKey = b

	return i
}

func (i *Item) SetField(index int, v Item) *Item {
	for len(i.Fields) <= index {
		i.Fields = append(i.Fields, Item{})
	}

	i.Fields[index] = v

	return i
}

func (i *Item) AddField(v Item) *Item {
	i.Fields = append(i.Fields, v)

	return i
}

func (i *Item) String() string {
	b, _ := json.MarshalIndent(i, "", "   ")

	return string(b)
}

func NewJSON(data []byte) (*Item, error) {
	var item Item

	err := json.Unmarshal(data, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
