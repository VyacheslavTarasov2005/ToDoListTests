package interfaces

import (
	"HITS_ToDoList_Tests/internal/application/enums"
	"HITS_ToDoList_Tests/internal/domain/models"
)

type TasksRepository interface {
	Add(task models.Task) error
	GetAll(sorting *enums.Sorting) ([]*models.Task, error)
}
