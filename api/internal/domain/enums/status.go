package enums

type Status string

const (
	Active    Status = "Active"
	Completed Status = "Completed"
	Overdue   Status = "Overdue"
	Late      Status = "Late"
)
