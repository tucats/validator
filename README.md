# validator v0.1.11

This is a JSON validator package. It allows tags to be added to structure definitions, and those structures are
then passed to a Define() operation which creates a map of the validation structure definitions. Subsequently, JSON
strings can be validated against the structure definitions to report if they are conformant or not.

This is intended to help catch misspelled fields, missing required fields, and invalid field values.

## Using the `validate` tag

Use the Go tag `validate` to identify what validations will be done on the JSON representation of the structure.
The tag is followed by a quoted string containing comma-separate validation operations. These operations are:

| Operation | Operand | Description |
| --------- | ------- | ----------- |
| required | | If specified, this field _must_ appear in the JSON |
| min | any | The minimum int or float value allowed for this field |
| max | any | The maximum int or float value allowed for this field |
| minlen | integer | The minimum length of a string value, or smallest allowed array size |
| maxlen | integer | The maximum length of a string value. or largest allowed array size |
| enum | strings | A list of strings separated by vertical bars enumerating the allowed field values |
| list | | The string value can be a list, each of which must match the enum list |
| matchcase | | The enumerated values must match case to match the field value |
| key | (items) | Specify limits on a map key value (which is always a string) |
| value | (items) | Specify rules on a value for an array or map |

You can separate enumerated values using commas rather than vertical bars by enclosing the
list of enumerated values in parenthesis. That is, `enum=red|green|blue` is the same as
specifying `enum=(red,green,blue)`. Note that leading and trailing spaces in enumerated values
are ignored.

some operations cannot be performed on all data types. For example, `min` and `max` can be
used with a time.Time value to compare the time provided in the JSON to specific time values. However,
these are not applicable to fields containing maps. For an array, the validations apply to the array
itself (such as minimum or maximum length) and the `value=` clause defines the validation rules for
the individual values in the array. Note that you can use parentheses to delimit lists in the `value=`
clause, such as

```go
type Foo struct {
    Colors []string `validate:"minlen=1,value=(enum=red,green,blue)"`
}
```

In this example, the JSON representation of a `Foo` structure must include an array with at least
one element (`minlen`) but each value in the array must also conform to an enumerated list allowing
only the values `red`, `green`, and `blue`.

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
builds a data structure that defines how the validation is to be performed. For example,

```go
    employeeValidator, err := validator.New(&Employees{})
```

The error return can indicate invalid tags or unsupportable data types for validation in the
data structure.

The `New()` function scans the object (a structure, in this case) and creates validation rules for any items specified with structure tags. It can support numeric, string, and boolean fields. It can handle the cases of uuid.UUID, time.Time, and time.Duration when expressed as a string. And it can handle arrays, pointers, and map types.

For integer value types, a default min and max value is automatically created based on the
size of the integer type. So a structure of Go type `uint8` will automatically have a minimum value of 0 and a maximum value of 255. Similarly, a structure field of type `float32` will have a size
range based on a 32-bit floating point value. No such checks are done for `float64` or `int64`
data types.

## Validating a JSON string

To validate a JSON string to see if it contains a valid representation of the object, use
the `Validate()` method for the validator object previously created.

```go
    // Read the JSON file and validate it's contents.
    b, err := io.ReadFile("my.json")
    if err == nil {
        err = employee.Validate(string(b))
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
validator for non-standard sizes of integer or float values, such as Int16 or
Float32.

Once you have created the validator, you can set attributes on it. For example,
this code sets a minimum and maximum value for the validator.

```go
    i := validator.NewType(validator.TypeInteger).SetMinValue(1).SetMaxValue(10)

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
| -------- | ----------- |
| SetMinValue(v) | Set the minimum allowed numeric value |
| SetMAxValue(v) | Set the maximum allowed numeric value |
| SetMinLen(i) | Set the minimum string or array length |
| SetMaxLen(i) | Set the maximum string or array length |
| SetEnum(v...) | Set the allowed values for integer or string values |
| SetField(i, v) | Set structure field `i` to validator `v` |
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
