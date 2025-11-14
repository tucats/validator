package tests

import (
	"reflect"
	"testing"

	"github.com/tucats/validator"
)

func TestNewJSON(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		data string
		want *validator.Item
		err  error
	}{
		{
			name: "Valid JSON for a string object",
			data: `{
					"type": "string"
					}`,
			want: &validator.Item{
				ItemType: validator.TypeString,
			},
		},
		{
			name: "Invalid JSON, misspelled key",
			data: `{
					"types": "strings"
					}`,
			err: validator.ErrInvalidValidator.Context("types").Value("invalid field name"),
		},
		{
			name: "Invalid JSON, bad type",
			data: `{
					"type": "strings"
					}`,
			err: validator.ErrInvalidValidator.Context("type").Value("missing or invalid type"),
		},
		{
			name: "Invalid JSON, min length without hasMinLength",
			data: `{
					"type": "string",
					"min_length": 5
					}`,
			err: validator.ErrInvalidValidator.Context("HasMinLength").Value("non-zero minLength without hasMinLength"),
		},
		{
			name: "Invalid JSON, nested validator has illegal type",
			data: `{
					"type": "pointer",
					"base_type": {
					    "type": "strings"
		                }
					}`,
			err: validator.ErrInvalidValidator.Context("type").Value("missing or invalid type"),
		},
		{
			name: "Invalid JSON, nested validator has illegal key",
			data: `{
					"type": "pointer",
					"base_type": {
					    "typ": "string"
		                }
					}`,
			err: validator.ErrInvalidValidator.Context("typ").Value("invalid field name"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := validator.NewJSON([]byte(tt.data))
			m1 := ""

			if gotErr != nil {
				m1 = gotErr.Error()
			}

			m2 := ""
			if tt.err != nil {
				m2 = tt.err.Error()
			}

			if m1 != m2 {
				t.Errorf("NewJSON() unexpected error: %v", gotErr)

				return
			}

			if gotErr == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Export(t *testing.T) {
	// Create a validator for the Employees structure.
	i1, err := validator.New(&Employees{})
	if err != nil {
		t.Errorf("Unexpected error creating validator: %s", err)

		return
	}

	// Export it's representation to JSON text.
	text := i1.String()

	// Use that text to recreate the validator.
	i2, err := validator.NewJSON([]byte(text))
	if err != nil {
		t.Errorf("Unexpected error creating validator from JSON: %s", err)

		return
	}

	// See if the two validators are identical.
	if !reflect.DeepEqual(i1, i2) {
		t.Error("Unexpected difference between original and recreated validators\n\nOriginal:\n", i1.String(), "\n\nRecreated:\n", i2.String())
	}
}
