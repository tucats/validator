package tests

import (
	"testing"

	"github.com/tucats/validator"
)

// Structure to validate.
type ListObject struct {
	Colors string `json:"colors" validate:"list,enum=red|green|blue,required,minlength=1,maxlength=3"`
	States string `json:"states" validate:"list,matchcase,enum=CA|NC|VT|TX,required,minlength=1,maxlength=4"`
}

func Test_Lists(t *testing.T) {
	type TestITem struct {
		name     string
		object   any
		jsonText string
		expected error
	}

	tests := []TestITem{
		{
			"Valid list",
			&ListObject{},
			`{
			    "colors": "red,blue",
				"states": "CA,VT"
			}`,
			nil,
		},
		{
			"invalid list, bad color",
			&ListObject{},
			`{
			    "colors": "red,pink",
				"states": "CA,VT"
			}`,
			validator.ErrInvalidEnumeratedValue.Context("colors").Value("pink").Expected([]string{"red", "green", "blue"}),
		},
		{
			"invalid list, too many colors",
			&ListObject{},
			`{
			    "colors": "red,blue,green,red",
				"states": "CA,VT"
			}`,
			validator.ErrValueLengthOutOfRange.Context("colors").Value("red,blue,green,red"),
		},
		{
			"invalid list, items not match case",
			&ListObject{},
			`{
			    "colors": "red,blue,green",
				"states": "ca,vt"
			}`,
			validator.ErrInvalidEnumeratedValue.Context("states").Value("ca").Expected([]string{"CA", "NC", "VT", "TX"}),
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
