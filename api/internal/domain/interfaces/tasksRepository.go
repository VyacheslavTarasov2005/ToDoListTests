package interfaces

import "HITS_ToDoList_Tests/internal/domain/models"

type TasksRepository interface {
	Add(task models.Task) error
}
