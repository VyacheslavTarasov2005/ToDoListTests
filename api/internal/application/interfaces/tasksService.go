package interfaces

import (
	"HITS_ToDoList_Tests/internal/domain/enums"
	"HITS_ToDoList_Tests/internal/domain/models"
	"time"
)

type TasksService interface {
	CreateTask(name string, description *string, deadline *time.Time, priority *enums.Priority) (*models.Task, error)
}
