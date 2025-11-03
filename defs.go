package validator

// Structures to validate.

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

type Person struct {
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Address Address `json:"address"`
}

type Employees struct {
	Departments string   `json:"department"`
	Staff       []Person `json:"staff"`
}
