package validators

import (
	"HITS_ToDoList_Tests/internal/application/errors"
	"time"
)

func ValidateTask(name string, deadline *time.Time) error {
	err := errors.ApplicationError{
		StatusCode: 400,
		Code:       "ValidationFailed",
		Errors:     map[string]string{},
	}

	if len(name) < 4 {
		err.Errors["name"] = "Name is required"
	}

	if deadline != nil && !deadline.After(time.Now()) {
		err.Errors["deadline"] = "Deadline must be in the future"
	}

	if len(err.Errors) > 0 {
		return err
	}

	return nil
}
