package tests

import (
	"reflect"
	"testing"

	"github.com/tucats/validator"
)

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
