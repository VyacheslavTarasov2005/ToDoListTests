package enums

type Sorting string

const (
	CreateAsc    Sorting = "CreateAsc"
	CreateDesc   Sorting = "CreateDesc"
	PriorityAsc  Sorting = "PriorityAsc"
	PriorityDesc Sorting = "PriorityDesc"
	DeadlineAsc  Sorting = "DeadlineAsc"
	DeadlineDesc Sorting = "DeadlineDesc"
)
