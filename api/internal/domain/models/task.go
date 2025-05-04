package models

import (
	"HITS_ToDoList_Tests/internal/domain/enums"
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	ChangedAt   *time.Time
	Name        string
	Description *string
	Deadline    *time.Time
	Status      enums.Status
	Priority    enums.Priority
}

func NewTask(name string, description *string, deadline *time.Time, status *enums.Status,
	priority *enums.Priority) *Task {
	task := &Task{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		Name:        name,
		Description: description,
		Deadline:    deadline,
	}

	if status == nil {
		task.Status = enums.Active
	} else {
		task.Status = *status
	}

	if priority == nil {
		task.Priority = enums.Medium
	} else {
		task.Priority = *priority
	}

	return task
}
