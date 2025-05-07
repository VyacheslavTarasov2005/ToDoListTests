package enums

import "fmt"

type Sorting string

const (
	CreateAsc    = "CreateAsc"
	CreateDesc   = "CreateDesc"
	PriorityAsc  = "PriorityAsc"
	PriorityDesc = "PriorityDesc"
	DeadlineAsc  = "DeadlineAsc"
	DeadlineDesc = "DeadlineDesc"
)

func ValidateSorting(s Sorting) error {
	switch s {
	case CreateAsc, CreateDesc, PriorityAsc, PriorityDesc, DeadlineAsc, DeadlineDesc:
		return nil
	default:
		return fmt.Errorf("invalid Sorting: %q", s)
	}
}
