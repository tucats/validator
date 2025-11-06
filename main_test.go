package validator

import (
	"testing"
)

func Test_AddressStruct(t *testing.T) {
	type TestITem struct {
		name     string
		object   any
		jsonText string
		expected error
	}

	tests := []TestITem{
		{
			"valid Employees",
			&Employees{},
			`{
			    "department": "Space Research",
				"division": "Engineering",
			    "staff": [
					{
						"name": "John Doe",
						"age": 35,
						"address":{
							"street": "123 Main St",
							"city": "New York"
						}
					},
					{
						"name": "Sue Smith",
						"age": 52,
						"address":{
							"street": "155 Oak Ave",
							"city": "New York"
						}
					}
			    ]
            }`,
			nil,
		},
		{
			"invalid Employees, bad division enum value",
			&Employees{},
			`{
			    "department": "Space Research",
				"division": "Science",
			    "staff": [
					{
						"name": "John Doe",
						"age": 35,
						"address":{
							"street": "123 Main St",
							"city": "New York"
						}
					},
					{
						"name": "Sue Smith",
						"age": 52,
						"address":{
							"street": "155 Oak Ave",
							"city": "New York"
						}
					}
			    ]
            }`,
			ErrInvalidEnumeratedValue.Context("division"),
		},
		{
			"invalid Employees, age out of range",
			&Employees{},
			`{
			    "department": "Space Research",
				"division": "Engineering",
			    "staff": [
					{
						"name": "John Doe",
						"age": 75,
						"address":{
							"street": "123 Main St",
							"city": "New York"
						}
					},
					{
						"name": "Sue Smith",
						"age": 52,
						"address":{
							"street": "155 Oak Ave",
							"city": "New York"
						}
					}
			    ]
            }`,
			ErrValueOutOfRange.Context("age"),
		},
		{
			"invalid Employees, empty staff array",
			&Employees{},
			`{
			    "department": "Space Research",
				"division": "Engineering",
			    "staff": []
            }`,
			ErrArrayLengthOutOfRange.Context("staff"),
		},
		{
			"valid JSON for address",
			&Address{},
			`{
				"street": "123 Main St",
				"city": "New York"
			}`,
			nil,
		},
		{
			"street string is too short",
			&Address{},
			`{
				"street": "",
				"city": "New York"
			}`,
			ErrValueLengthOutOfRange.Context("street"),
		},
		{
			"city field not present",
			&Address{},
			`{
				"street": "123 Main St"
			}`,
			ErrRequired.Context("city"),
		},
		{
			"valid Person",
			&Person{},
			`{
				"name": "John Doe",
				"age": 35,
				"address":{
					"street": "123 Main St",
					"city": "New York"
				}
			}`,
			nil,
		},
		{
			"invalid Person, age out of range",
			&Person{},
			`{
				"name": "John Doe",
				"age": 15,
				"address":{
					"street": "123 Main St",
					"city": "New York"
				}
			}`,
			ErrValueOutOfRange.Context("age"),
		},
		{
			"invalid Person, missing city field",
			&Person{},
			`{
				"name": "John Doe",
				"age": 42,
				"address":{
					"street": "123 Main St"
				}
			}`,
			ErrRequired.Context("city"),
		},
	}

	for _, test := range tests {
		var (
			msg1 string
			msg2 string
		)

		// Define the structures we want to validate.
		item, err := New(test.object)
		if err != nil {
			t.Fatal("Failed to define Address structure:", err)
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
