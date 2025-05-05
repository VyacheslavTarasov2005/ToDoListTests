package services

import (
	appEnums "HITS_ToDoList_Tests/internal/application/enums"
	"HITS_ToDoList_Tests/internal/application/errors"
	appInterfaces "HITS_ToDoList_Tests/internal/application/interfaces"
	"HITS_ToDoList_Tests/internal/application/validators"
	"HITS_ToDoList_Tests/internal/domain/enums"
	domainInterfaces "HITS_ToDoList_Tests/internal/domain/interfaces"
	"HITS_ToDoList_Tests/internal/domain/models"
	"github.com/google/uuid"
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

func (service *TasksServiceImpl) DeleteTask(taskID uuid.UUID) error {
	task, err := service.tasksRepository.GetByID(taskID)
	if err != nil {
		return err
	}

	if task == nil {
		return errors.ApplicationError{
			StatusCode: 404,
			Code:       "NotFound",
			Errors:     map[string]string{"message": "Task not found"},
		}
	}

	err = service.tasksRepository.DeleteByID(taskID)
	if err != nil {
		return err
	}

	return nil
}
