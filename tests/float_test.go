package tests

import (
	"testing"

	"github.com/tucats/validator"
)

// Structure to validate.
type FloatObject struct {
	Item32 float32 `json:"item32"`
	Item64 float64 `json:"item64"`
}

func Test_Floats(t *testing.T) {
	type TestITem struct {
		name     string
		object   any
		jsonText string
		expected error
	}

	tests := []TestITem{
		{
			"Valid float32",
			&FloatObject{},
			`{
			    "item32": 500.25
			}`,
			nil,
		},
		{
			"Invalid float32, too large",
			&FloatObject{},
			`{
			    "item32": 1.0e305
			}`,
			validator.ErrValueOutOfRange.Context("item32").Value(1.0e305),
		},
	}

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
