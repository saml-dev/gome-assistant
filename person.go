package gomeassistant

import "fmt"

type Person struct {
	// Personal details
	name, address, pin string
	// Job details
	workAddress, company, position string
	salary                         int
}

// PersonBuilder struct
type PersonBuilder struct {
	person *Person
}

// PersonAddressBuilder facet of PersonBuilder
type PersonAddressBuilder struct {
	PersonBuilder
}

// PersonJobBuilder facet of PersonBuilder
type PersonJobBuilder struct {
	PersonBuilder
}

// NewPersonBuilder constructor for PersonBuilder
func NewPersonBuilder() *PersonBuilder {
	return &PersonBuilder{person: &Person{}}
}

// Lives chains to type *PersonBuilder and returns a *PersonAddressBuilder
func (b *PersonBuilder) Lives() *PersonAddressBuilder {
	return &PersonAddressBuilder{*b}
}

// Works chains to type *PersonBuilder and returns a *PersonJobBuilder
func (b *PersonBuilder) Works() *PersonJobBuilder {
	return &PersonJobBuilder{*b}
}

// At adds address to person
func (a *PersonAddressBuilder) At(address string) *PersonAddressBuilder {
	a.person.address = address
	return a
}

// WithPostalCode adds postal code to person
func (a *PersonAddressBuilder) WithPostalCode(pin string) *PersonAddressBuilder {
	a.person.pin = pin
	return a
}

// As adds position to person
func (j *PersonJobBuilder) As(position string) *PersonJobBuilder {
	j.person.position = position
	return j
}

// For adds company to person
func (j *PersonJobBuilder) For(company string) *PersonJobBuilder {
	j.person.company = company
	return j
}

// In adds company address to person
func (j *PersonJobBuilder) In(companyAddress string) *PersonJobBuilder {
	j.person.workAddress = companyAddress
	return j
}

// WithSalary adds salary to person
func (j *PersonJobBuilder) WithSalary(salary int) *PersonJobBuilder {
	j.person.salary = salary
	return j
}

// Build builds a person from PersonBuilder
func (b *PersonBuilder) Build() *Person {
	return b.person
}

// RunBuilderFacet example
func RunBuilderFacet() {
	pb := NewPersonBuilder()
	pb.
		Lives().
		At("Bangalore").
		WithPostalCode("560102").
		Works().
		As("Software Engineer").
		For("IBM").
		In("Bangalore").
		WithSalary(150000)

	person := pb.Build()

	fmt.Println(person)
}
