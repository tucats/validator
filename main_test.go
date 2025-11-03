package validator

import (
	"testing"
)

func Test_Driver(t *testing.T) {
	type TestITem struct {
		name     string
		jsonText string
		expected bool
	}

	tests := []TestITem{
		{
			"Valid JSON for address",
			`{
				"street": "123 Main St",
				"city": "New York",
			}`,
			true,
		},
	}

	// Define the structures we want to validate.
	err := Define("Address", &Address{})
	if err != nil {
		t.Fatal("Failed to define Address structure:", err)
	}

	for _, test := range tests {
		err := Validate(test.name, test.jsonText)
		expected := (err == nil) == test.expected

		if !expected {
			t.Fatalf("JSON file %s validity: %v\n", test.jsonText, err)
		}
	}
}
