package tests

import (
	"testing"

	"github.com/tucats/validator"
)

// Test creating a validator item and adding attributes to it programmatically,
// rather than by parsing the information in a struct tag.
func Test_Bespoke(t *testing.T) {
	i, err := validator.New(0)
	if err != nil {
		t.Fatalf("Unexpected error creating validator: %v", err)
	}

	// Name the validator, and set minium and maximum allowed values.
	i = i.SetName("foo").SetMinValue(10).SetMaxValue(100)

	err = i.Validate("3")
	if err.Error() != validator.ErrValueOutOfRange.Context("foo").Value(3).Error() {
		t.Fatalf("Unexpected error %v", err)
	}

	err = i.Validate("15")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// Further restrict the validation to require one of these enumerated values.
	i.SetEnums(20, 30, 40)

	e2 := validator.ErrInvalidEnumeratedValue.Context("foo").Value(15).Expected(20, 30, 40)
	err = i.Validate("15")

	if err.Error() != e2.Error() {
		t.Fatalf("Unexpected error %v, got %v", e2, err)
	}
}
