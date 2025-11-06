package validator

// Structures to validate.

type Address struct {
	Street string `json:"street" validate:"required,minlength=1,maxlength=100"`
	City   string `json:"city"   validate:"required,minlength=1,maxlength=100"`
}

type Person struct {
	Name    string  `json:"name"    validate:"required,minlength=1,maxlength=100"`
	Age     int     `json:"age"     validate:"required,min=18,max=65"`
	Address Address `json:"address" validate:"required"`
}

type Employees struct {
	Department string   `json:"department" validate:"required"`
	Division   string   `json:"division"   validate:"required,enum=HR|Finance|Marketing|Engineering"`
	Staff      []Person `json:"staff"      validate:"minlen=1"`
}
