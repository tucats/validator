package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/tucats/validator"
)

func Test_Integration(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		src  string
		json string
		err  error
	}{
		{
			name: "simple integer",
			src:  `int: minvalue=1, maxvalue=10`,
			json: `9`,
			err:  nil,
		},
		{
			name: "simple integer, out of range",
			src:  `int: minvalue=1, maxvalue=10`,
			json: `42`,
			err:  validator.ErrValueOutOfRange.Value(42),
		},
		{
			name: "array of integers",
			src:  `[]int: minvalue=1, maxvalue=10`,
			json: `[9, 3, 1]`,
			err:  nil,
		},
		{
			name: "array of integers with value out of range",
			src:  `[]int: base=(minvalue=1, maxvalue=10)`,
			json: `[9, 13, 1]`,
			err:  validator.ErrValueOutOfRange.Value(13),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := tt.src
			if strings.HasPrefix(text, "@") {
				b, err := os.ReadFile(tt.src[1:])
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}

				text = string(b)
			}

			v, gotErr := validator.Compile(text)
			if gotErr != nil {
				t.Fatalf("failed to compile JSON: %v", gotErr)
			}

			source := tt.json
			if strings.HasPrefix(source, "@") {
				b, err := os.ReadFile(tt.json[1:])
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}

				source = string(b)
			}

			gotErr = v.Validate(source)

			var m1, m2 string

			if gotErr != nil {
				m1 = gotErr.Error()
			}

			if tt.err != nil {
				m2 = tt.err.Error()
			}

			if m1 != m2 {
				t.Errorf("Compile(%s) unexpected error\n  wanted: %v\n. got:    %v", tt.name, tt.err, gotErr)

				return
			}
		})
	}
}
