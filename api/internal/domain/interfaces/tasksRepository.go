package interfaces

import (
	"HITS_ToDoList_Tests/internal/application/enums"
	"HITS_ToDoList_Tests/internal/domain/models"
	"github.com/google/uuid"
)

type TasksRepository interface {
	Add(task models.Task) error
	GetAll(sorting *enums.Sorting) ([]*models.Task, error)
	GetByID(id uuid.UUID) (*models.Task, error)
	DeleteByID(taskID uuid.UUID) error
}
