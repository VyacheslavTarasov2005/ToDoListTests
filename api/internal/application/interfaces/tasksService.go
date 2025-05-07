package interfaces

import (
	appEnums "HITS_ToDoList_Tests/internal/application/enums"
	"HITS_ToDoList_Tests/internal/domain/enums"
	"HITS_ToDoList_Tests/internal/domain/models"
	"github.com/google/uuid"
	"time"
)

type TasksService interface {
	CreateTask(name string, description *string, deadline *time.Time, priority *enums.Priority) (*models.Task, error)
	GetAllTasks(sorting *appEnums.Sorting) ([]*models.Task, error)
	DeleteTask(taskID uuid.UUID) error
	UpdateTask(taskID uuid.UUID, name string, description *string, deadline *time.Time,
		priority *enums.Priority) (*models.Task, error)
	ToggleTaskStatus(taskID uuid.UUID, isDone bool) (*models.Task, error)
	UpdateTaskStatuses()
}
