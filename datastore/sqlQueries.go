package datastore

import (
	"database/sql/driver"
	"fmt"
	"github.com/ellis90/assessment-bg/entity"
)

type statusWrapper entity.Status

func (s statusWrapper) Value() (driver.Value, error) {
	switch entity.Status(s) {
	case entity.Inactive:
		return "I", nil
	case entity.Active:
		return "A", nil
	case entity.Terminated:
		return "T", nil
	default:
		return nil, fmt.Errorf("invalid Status %d", s)
	}
}

// Scan implements database/sql/driver.Scanner
func (s *statusWrapper) Scan(in any) error {
	switch in.(string) {
	case "I":
		*s = entity.Inactive
		return nil
	case "A":
		*s = entity.Active
		return nil
	case "T":
		*s = entity.Terminated
		return nil
	default:
		return fmt.Errorf("invalid Status: %q", in.(string))
	}
}
