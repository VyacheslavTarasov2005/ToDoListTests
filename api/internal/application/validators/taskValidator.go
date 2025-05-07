package validators

import (
	"HITS_ToDoList_Tests/internal/application/errors"
	"HITS_ToDoList_Tests/internal/domain/enums"
	"time"
)

func ValidateTask(name string, deadline *time.Time, priority *enums.Priority) error {
	err := errors.ApplicationError{
		StatusCode: 400,
		Code:       "ValidationFailed",
		Errors:     map[string]string{},
	}

	if name == "" {
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
