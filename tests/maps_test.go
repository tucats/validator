package tests

import (
	"testing"

	"github.com/tucats/validator"
)

// Structure to validate.
type MapStrings struct {
	Items map[string]string `json:"items" validate:"required, key=(key1,key2), value=(enum=(value1, value2))"`
}

type MapInts struct {
	Items map[string]int `json:"items" validate:"required,enum=key1|key2"`
}

type MapStringArray struct {
	Items map[string][]string `json:"items" validate:"required,enum=key1|key2,base=(enum=value1|value2|value3|value4)"`
}

func Test_Maps(t *testing.T) {
	type TestITem struct {
		name     string
		object   any
		jsonText string
		expected error
	}

	tests := []TestITem{
		{
			"Valid map[string]string",
			&MapStrings{},
			`{
			    "items": {
				    "key1": "value1",
                    "key2": "value2"
				}
			}`,
			nil,
		},
		{
			"Valid map[string]int",
			&MapInts{},
			`{
			    "items": {
				    "key1": 55,
                    "key2": 67
				}
			}`,
			nil,
		},
		{
			"Valid map[string][]string]",
			&MapStringArray{},
			`{
			    "items": {
				    "key1": ["value1","value2"],
                    "key2": ["value3","value4"]
				}
			}`,
			nil,
		},
		{
			"invalid map key value",
			&MapInts{},
			`{
			    "items": {
				    "key3": 55,
                    "key2": 67
				}
			}`,
			validator.ErrInvalidEnumeratedValue.Context("items").Value("key3").Expected([]string{"key1", "key2"}),
		},
		{
			"Invalid map[string]string, bad value",
			&MapStrings{},
			`{
			    "items": {
				    "key1": "value3",
                    "key2": "value2"
				}
			}`,
			validator.ErrInvalidEnumeratedValue.Value("value3").Expected([]string{"value1", "value2"}),
		},
		{
			"Invalid map[string][]string], bad value array member",
			&MapStringArray{},
			`{
			    "items": {
				    "key1": ["value1","value5"],
                    "key2": ["value3","value4"]
				}
			}`,
			validator.ErrInvalidEnumeratedValue.Value("value5").Expected([]string{"value1", "value2", "value3", "value4"}),
		},
		{
			"invalid map, missing required field",
			&MapStrings{},
			`{
			}`,
			validator.ErrRequired.Value("items"),
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
