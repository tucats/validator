package tests

import (
	"testing"

	"github.com/tucats/validator"
)

// Structure to validate.
type IntObject struct {
	Item8   int8   `json:"item8"`
	Item16  int16  `json:"item16"`
	Item32  int32  `json:"item32"`
	ItemU8  uint8  `json:"itemU8"`
	ItemU16 uint16 `json:"itemU16"`
	ItemU32 uint32 `json:"itemU32"`
}

func Test_Ints(t *testing.T) {
	type TestITem struct {
		name     string
		object   any
		jsonText string
		expected error
	}

	tests := []TestITem{
		{
			"Valid uint8",
			&IntObject{},
			`{
			    "itemU8": 50
			}`,
			nil,
		},
		{
			"Invalid uint8, too large",
			&IntObject{},
			`{
			    "itemU8": 256
			}`,
			validator.ErrValueOutOfRange.Context("itemU8").Value(256),
		},
		{
			"Invalid uint8, too small",
			&IntObject{},
			`{
			    "itemU8": -5
			}`,
			validator.ErrValueOutOfRange.Context("itemU8").Value(-5),
		},
		{
			"Valid uint16",
			&IntObject{},
			`{
			    "itemU16": 5000
			}`,
			nil,
		},
		{
			"Invalid uint16, too large",
			&IntObject{},
			`{
			    "itemU16": 66000
			}`,
			validator.ErrValueOutOfRange.Context("itemU16").Value(66000),
		},
		{
			"Invalid uint16, too small",
			&IntObject{},
			`{
			    "itemU16": -1000
			}`,
			validator.ErrValueOutOfRange.Context("itemU16").Value(-1000),
		},
		{
			"Valid uint32",
			&IntObject{},
			`{
			    "itemU32": 66000
			}`,
			nil,
		},
		{
			"Invalid int32, too large",
			&IntObject{},
			`{
			    "itemU32": 50000000000
			}`,
			validator.ErrValueOutOfRange.Context("itemU32").Value("50000000000"),
		},
		{
			"Valid int8",
			&IntObject{},
			`{
			    "item8": 50
			}`,
			nil,
		},
		{
			"Invalid int8, too large",
			&IntObject{},
			`{
			    "item8": 1000
			}`,
			validator.ErrValueOutOfRange.Context("item8").Value(1000),
		},
		{
			"Valid int16",
			&IntObject{},
			`{
			    "item16": 5000
			}`,
			nil,
		},
		{
			"Invalid int16, too large",
			&IntObject{},
			`{
			    "item16": 50000
			}`,
			validator.ErrValueOutOfRange.Context("item16").Value(50000),
		},
		{
			"Valid int32",
			&IntObject{},
			`{
			    "item32": 50000
			}`,
			nil,
		},
		{
			"Invalid int32, too large",
			&IntObject{},
			`{
			    "item32": 50000000000
			}`,
			validator.ErrValueOutOfRange.Context("item32").Value("50000000000"),
		}}

	for _, test := range tests {
		var (
			msg1 string
			msg2 string
		)

		// Define the structures we want to validate.
		item, err := validator.New(test.object)
		if err != nil {
			t.Fatal("Failed to define structure:", err)
		}

		err = item.Validate(test.jsonText)
		if err != nil {
			msg1 = err.Error()
		}

		if test.expected != nil {
			msg2 = test.expected.Error()
		}

		expected := (msg1 == msg2)

		if !expected {
			t.Fatalf("In \"%s\", unexpected result: %v\n", test.name, err)
		}
	}
}
