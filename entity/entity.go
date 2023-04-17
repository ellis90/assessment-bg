package entity

import "fmt"

const (
	Inactive = iota
	Active
	Terminated
)

type Status int

func (s Status) String() string {
	switch s {
	case Inactive:
		return "I"
	case Active:
		return "A"
	case Terminated:
		return "T"
	default:
		return fmt.Sprintf("invalid Status %d", s)
	}
}

// User entity represents every user in the domain
type User struct {
	ID         string `json:"id"`
	UserName   string `json:"userName" validate:"required"`
	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Department string `json:"department" validate:"required"`
	UserStatus Status `json:"userStatus" validate:"gte=0,lte=2"`
}
