package enums

type Status int

const (
	Active Status = iota
	Completed
	Overdue
	Late
)
