package validators

import (
	"HITS_ToDoList_Tests/internal/application/errors"
	"HITS_ToDoList_Tests/internal/domain/enums"
	"HITS_ToDoList_Tests/internal/domain/models"
	"time"
)

func ValidateTask(task models.Task) error {
	err := errors.ApplicationError{
		StatusCode: 400,
		Code:       "Validation",
		Errors:     map[string]string{},
	}

	if task.Name == "" {
		err.Errors["name"] = "Name is required"
	}

	if task.Deadline != nil && !task.Deadline.After(time.Now()) {
		err.Errors["deadline"] = "Deadline must be in the future"
	}

	validPriorities := map[enums.Priority]bool{
		enums.Low:      true,
		enums.Medium:   true,
		enums.High:     true,
		enums.Critical: true,
	}

	if !validPriorities[task.Priority] {
		err.Errors["priorities"] = "Priorities is required"
	}

	if len(err.Errors) > 0 {
		return err
	}

	return nil
}
