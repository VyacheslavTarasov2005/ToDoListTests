package services

import (
	appEnums "HITS_ToDoList_Tests/internal/application/enums"
	appInterfaces "HITS_ToDoList_Tests/internal/application/interfaces"
	"HITS_ToDoList_Tests/internal/application/validators"
	"HITS_ToDoList_Tests/internal/domain/enums"
	domainInterfaces "HITS_ToDoList_Tests/internal/domain/interfaces"
	"HITS_ToDoList_Tests/internal/domain/models"
	"time"
)

type TasksServiceImpl struct {
	tasksRepository domainInterfaces.TasksRepository
}

func NewTasksService(tasksRepository domainInterfaces.TasksRepository) appInterfaces.TasksService {
	return &TasksServiceImpl{tasksRepository: tasksRepository}
}

func (service *TasksServiceImpl) CreateTask(
	name string, description *string, deadline *time.Time, priority *enums.Priority) (*models.Task, error) {
	task := models.NewTask(name, description, deadline, nil, priority)
	if err := validators.ValidateTask(*task); err != nil {
		return nil, err
	}

	if err := service.tasksRepository.Add(*task); err != nil {
		return nil, err
	}

	return task, nil
}

func (service *TasksServiceImpl) GetAllTasks(sorting *appEnums.Sorting) ([]*models.Task, error) {
	tasks, err := service.tasksRepository.GetAll(sorting)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
