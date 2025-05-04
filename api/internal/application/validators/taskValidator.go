package validators

import (
	"HITS_ToDoList_Tests/internal/application/errors"
	"HITS_ToDoList_Tests/internal/domain/models"
	"time"
)

func ValidateTask(task models.Task) error {
	err := errors.ApplicationError{
		StatusCode: 400,
		Code:       "validation",
		Errors:     map[string]string{},
	}

	if task.Name == "" {
		err.Errors["name"] = "Name is required"
	}

	if task.Deadline != nil && task.Deadline.Before(time.Now()) {
		err.Errors["deadline"] = "Deadline can't be in the past"
	}

	if len(err.Errors) > 0 {
		return err
	}

	return nil
}
