package datastore

import (
	"errors"
	"fmt"
	"github.com/ellis90/assessment-bg/datastore/model"
)

const (
	usersSchema = "users"
)

var (
	ErrDeleteCustomer         = errors.New("failed to delete customer")
	ErrFailedToCreateCustomer = errors.New("failed to add customer")
	ErrUpdateCustomer         = errors.New("failed to update customer")
	ErrFetchCustomer          = errors.New("failed to fetch customer")
	errorMsg                  = fmt.Sprintf("%e : %w\n")
)

type UserRepository interface {
	Create(user model.Customer) (model.Customer, error)
	Update(user model.Customer) (model.Customer, error)
	Get() (model.Customers, error)
	Delete(id string) error
}
