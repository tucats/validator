package tests

import (
	"testing"

	"github.com/tucats/validator"
)

func Test_Recursion(t *testing.T) {
	// Create a structure that allows recursion
	type DejaVu struct {
		Name     string   `json:"name"     validate:"required,minlength=5,maxlength=100"`
		Children []DejaVu `json:"children" validate:"minlength=0"`
	}

	// Create a validator for the DejaVu structure.
	i, err := validator.New(&DejaVu{})
	if err != nil {
		t.Errorf("Unexpected error creating validator: %s", err)

		return
	}

	// Make a structure that contains a recursive reference.
	text := `{
		"name": "parent",
		"children": [
			{
				"name": "child1",
				"children": [
					{
						"name": "grandchild1",
						"children": []
					},
					{
						"name": "grandchild2",
						"children": []
					}
				]
			},
			{
				"name": "child2",
				"children": []
			}
		]
	}`

	// Validate the structure.
	err = i.Validate(text)
	if err != nil {
		t.Errorf("Unexpected error validating recursive structure: %s", err)
	}

	// Do it again with a structure with an invalid field value in the recursive reference.
	text = `{
		"name": "parent",
		"children": [
			{
				"name": "child1",
				"children": [
					{
						"name": "grandchild1",
						"children": []
					},
					{
						"name": "zrg",
						"children": []
					}
				]
			},
			{
				"name": "child2",
				"children": []
			}
		]
	}`

	// Validate the structure.
	err = i.Validate(text)
	msg1 := err.Error()
	msg2 := validator.ErrValueLengthOutOfRange.Context("name").Value("zrg").Error()

	if msg1 != msg2 {
		t.Errorf("Unexpected error validating recursive structure: %s", err)
	}
}
