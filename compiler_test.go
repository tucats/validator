package validator_test

import (
	"reflect"
	"testing"

	"github.com/tucats/validator"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		src  string
		want *validator.Item
		err  error
	}{
		{
			name: "simple integer type",
			src:  "int;",
			want: &validator.Item{
				ItemType: validator.TypeInt,
			},
			err: nil,
		},
		{
			name: "named integer type",
			src:  "age int;",
			want: &validator.Item{
				Name:     "age",
				ItemType: validator.TypeInt,
			},
			err: nil,
		},
		{
			name: "simple integer type with comment",
			src: `// Compiled validator for an integer
			int;
			`,
			want: &validator.Item{
				ItemType: validator.TypeInt,
			},
			err: nil,
		},
		{
			name: "invalid type",
			src:  "integer;",
			err:  validator.ErrUnsupportedType.Context("line 1, column 8").Value("integer"),
		},
		{
			name: "integer type with simple tag",
			src:  "int: required;",
			want: &validator.Item{
				ItemType: validator.TypeInt,
				Required: true,
			},
			err: nil,
		},
		{
			name: "integer type with longer tag",
			src:  "int: required, minvalue=1, maxvalue=10",
			want: &validator.Item{
				ItemType:    validator.TypeInt,
				Required:    true,
				HasMinValue: true,
				MinValue:    "1",
				HasMaxValue: true,
				MaxValue:    "10",
			},
			err: nil,
		},
		{
			name: "integer type with invalid tag",
			src:  "int: omit=true",
			err:  validator.ErrInvalidKeyword.Value("omit"),
		},
		{
			name: "object type with auto-generated line endings",
			src: `person {
			       age int: required, 
				            minvalue=18, 
							maxvalue=65
				   name string: required, 
				                minlength=1, 
								maxlength=101
				   }`,
			want: &validator.Item{
				ItemType: validator.TypeStruct,
				Name:     "person",
				Fields: []*validator.Item{
					{
						Name:        "age",
						ItemType:    validator.TypeInt,
						Required:    true,
						HasMinValue: true,
						MinValue:    "18",
						HasMaxValue: true,
						MaxValue:    "65",
					},
					{
						Name:         "name",
						ItemType:     validator.TypeString,
						Required:     true,
						HasMinLength: true,
						MinLength:    1,
						HasMaxLength: true,
						MaxLength:    101,
					},
				},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := validator.Compile(tt.src)

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

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Compile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateLineEndings(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		src  string
		want string
	}{
		{
			name: "no line ending needed",
			src:  "int;",
			want: "int;",
		},
		{
			name: "line ending needed",
			src:  "int",
			want: "int;",
		},
		{
			name: "line ending needed",
			src: `person {
age int
}`,
			want: `person {
age int;
};`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.UpdateLineEndings(tt.src)
			if got != tt.want {
				t.Errorf("UpdateLineEndings() = %v, want %v", got, tt.want)
			}
		})
	}
}
