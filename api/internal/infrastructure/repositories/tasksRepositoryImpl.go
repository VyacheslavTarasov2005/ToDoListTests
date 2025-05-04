package repositories

import (
	"HITS_ToDoList_Tests/internal/domain/interfaces"
	"HITS_ToDoList_Tests/internal/domain/models"
	"gorm.io/gorm"
)

type TasksRepositoryImpl struct {
	db *gorm.DB
}

func NewTasksRepository(db *gorm.DB) interfaces.TasksRepository {
	return &TasksRepositoryImpl{db: db}
}

func (repo *TasksRepositoryImpl) Add(task models.Task) error {
	return repo.db.Create(task).Error
}
