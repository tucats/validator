# validator

This is a JSON validator package. It allows tags to be added to structure definitions, and those structures are
then passed to a Define() operation which creates a map of the valid structure definitions. Subsequently, JSON
strings can be validated against the structure definitions to report if they are conformant or not.

This is intended to help catch misspelled fields, missing require fields, and invalid field values.

## Using the `validate` tag

Use the Go tag `validate` to identify what validations will be done on the JSON representation of the structure.
The tag is followed by a quoted string containing comma-separate validation operations. These operations are:

| Operation | Operand | Description |
|-----------|---------|-------------|
| required | | If specified, this field _must_ appear in the JSON |
| min | any | The minimum int or float value allowed for this field |
| max | any | The maximum int or float value allowed for this field |
| minlen | integer | The minimum length of a string value, or smallest allowed array size |
| maxlen | integer | The maximum length of a string value. or largest allowed array size |
| enum | strings | A list of strings separated by vertical bars enumerating the allowed field values |
| matchcase| | The enumerated values must match case to match the field value |

Note that some operations cannot be performed on all data types. For example, `min` and `max` can be
used with a time.Time value to compare the time provided in the JSON to specific time values. However,
these are not applicable to fields containing maps. By contrast, a map can only support the `enum`
operator to validate the key values in the map. There are no validations for the values of the map.

Here is an example of a set of structures that are to be used to process JSON data. The associated `json`
and `validate` tags indicate how the field names are handled by JSON and the additional validation to be
done.

```go


type Address struct {
    Street string `json:"street" validate:"required,minlength=1,maxlen=100"`
    City   string `json:"city"   validate:"required,minlength=1,maxlen=100"`
}

type Person struct {
    Name    string  `json:"name"    validate:"required,minlen=1,maxlen=100"`
    Age     int     `json:"age"     validate:"required,min=18,max=65"`
    Address Address `json:"address" validate:"required"`
}

type Employees struct {
    Department string   `json:"department" validate:"required"`
    Division   string   `json:"division"   validate:"required,enum=HR|Finance|Marketing|Engineering"`
    Staff      []Person `json:"staff"      validate:"minlen=1"`
}
```

## Creating a new Validator object

A validator object is created by passing an instance of the structure to the validator, which
builds a data structure that defines how the validation is to be performed.  For example,

```go
    employeeValidator, err := validator.New(&Employees{})
```

The error return can indicate invalid tags or unsupportable data types for validation in the
data structure.

## Validating a JSON string

To validate a JSON string to see if it contains a valid representation of the object, use
the `Validate()` method for the validator object previously created.

```go
    // Read the JSON file and validate it's contents.
    b, err := io.ReadFile("my.json")
    if err == nil {
        err = employeeValidator(string(b))
    }

    // JSON text does not violate any specified data requirements...

```

This (oversimplified) example shows validating the JSON data (which must be represented as
a string) to verify that the JSON representation contains the required fields, no misspelled
field names, and no invalid values. If the error return is nil, no errors where found.

Currently, the validator stops on the first error it finds and reports it.

## Programmatic Validators

Validators can be created programmatically rather than by reading tags from Go
structure definitions. This is useful if validating a single object that is
not part of a structure, for example. Create the validator as usual, by passing
in an example of the data type to be validated.

```go
    i := validator.NewType(validator.TypeInteger)
```

In this example, a validator is created for an integer value (due to the use
of the value `0` in the call to `New()` to create the validator. Note that not
all Go types are supported by the validator. For example, you cannot create a
validator for non-standard integer or float values, such as Int15 or Float32.

Once you have created the validator, you can set attributes on it. For example,
this code sets a minimum and maximum value for the validator.

```go
    i := validator.NewType(validator.TypeInteger).SetMaxValue(1).SetMaxValue(10)

    err := i.Validate(15)
    if err == nil {
        return err
    }
    ...
```

In this example, the validator is created for an integer value, and has minimum
and maximum values set. A test to see if the value `15` is valid is performed.
Because the validator requires a value from 1..10, the validator will return
an error.

The following modifier functions are used to set the characteristics of a
validator.

| Function | Description |
|----------|-------------|
| SetMinValue(v)   | Set the minimum allowed numeric value |
| SetMAxValue(v)   | Set the maximum allowed numeric value |
| SetMinLen(i)     | Set the minimum string or array length |
| SetMaxLen(i)     | Set the maximum string or array length |
| SetEnum(v...)    | Set the allowed values for integer or string values |
| SetField(i, v)   | Set structure field `i` to validator `v` |
| AddField(v) | Add a new structure field to the validator |
| SetMatchCase(b) | Indicate if enumerated strings must match case |
| SetForeignKey(b) | Indicate if undeclared field names are permitted |

## Import and Export

A validator can be converted to a string representation as a JSON object using
the `String()` function.

A validator can be created by reading a JSON string (typically that was
created by the `String()` function) using the `NewFromJSON()` function.

These functions make it easier for code to be written to allow externally-
created validator definitions in JSON to be integrated into the program.
