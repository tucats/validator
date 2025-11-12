package tests

import (
	"testing"
	"time"

	"github.com/tucats/validator"
)

// Structure to validate.
type DurationObject struct {
	Wait time.Duration `json:"wait" validate:"required,min=1s,max=1m"`
}

func Test_Duration(t *testing.T) {
	type TestITem struct {
		name     string
		object   any
		jsonText string
		expected error
	}

	tests := []TestITem{
		{
			"Valid time",
			&DurationObject{},
			`{
			    "wait": "4500ms"
			}`,
			nil,
		},
		{
			"invalid duration format",
			&DurationObject{},
			`{
			    "wait": "Yesterday"
			}`,
			validator.ErrInvalidData.Context("wait").Value("Yesterday"),
		},
		{
			"wait too short",
			&DurationObject{},
			`{
			    "wait": "15ms"
			}`,
			validator.ErrValueOutOfRange.Context("wait").Value("15ms"),
		},
		{
			"wait too long",
			&DurationObject{},
			`{
			    "wait": "2h"
			}`,
			validator.ErrValueOutOfRange.Context("wait").Value("2h"),
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
