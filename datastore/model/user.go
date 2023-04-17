package model

import (
	"errors"
	"fmt"
	"github.com/ellis90/assessment-bg/entity"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidPerson = errors.New("a customer/user has a missing field")
	// use a single instance of Validate, it caches struct info
	validate *validator.Validate
)

type Customer struct {
	person *entity.User
}

type ExportCustomer struct {
	User entity.User
}

type Customers []Customer
type ExportCustomers []ExportCustomer

func NewCustomer(user *entity.User) (Customer, error) {
	validate = validator.New()
	if err := validateStruct(user); err != nil {
		return Customer{}, err
	}
	if err := validateVariable(user.Email); err != nil {
		return Customer{}, err
	}
	return Customer{person: user}, nil
}

func AddCustomer(user *entity.User) Customer {
	return Customer{person: user}
}

func (c *Customer) SetID(id string) {
	c.person.ID = id
}
func (c *Customer) GetID() string {
	return c.person.ID
}

func (c *Customer) GetFirstName() string {
	return c.person.FirstName
}
func (c *Customer) GetLastName() string {
	return c.person.LastName
}
func (c *Customer) GetUserName() string {
	return c.person.UserName
}

func (c *Customer) GetEmail() string {
	return c.person.Email
}

func (c *Customer) GetDepartment() string {
	return c.person.Department
}

func (c *Customer) GetUserStatus() string {
	return c.person.UserStatus.String()
}

func (c *Customer) GetExportedCustomer() ExportCustomer {
	return ExportCustomer{
		User: *c.person,
	}
}

func (cs Customers) GetExportedCustomers() ExportCustomers {
	var expc ExportCustomers
	for _, c := range cs {
		expc = append(expc, c.GetExportedCustomer())
	}
	return expc
}

func validateStruct(user *entity.User) error {
	// returns nil or ValidationErrors ( []FieldError )
	err := validate.Struct(user)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		var fieldErr error
		for _, err := range err.(validator.ValidationErrors) {
			fieldErr = fmt.Errorf("nameSpace: %s, field: %s\n", err.Namespace(), err.Field())
		}

		// from here you can create your own error messages in whatever language you wish
		return fmt.Errorf("%e: %w", fieldErr, err)
	}
	return nil
}

func validateVariable(email string) error {
	errs := validate.Var(email, "required,email")
	return errs
}
