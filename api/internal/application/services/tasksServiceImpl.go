package services

import (
	appEnums "HITS_ToDoList_Tests/internal/application/enums"
	"HITS_ToDoList_Tests/internal/application/errors"
	appInterfaces "HITS_ToDoList_Tests/internal/application/interfaces"
	"HITS_ToDoList_Tests/internal/application/validators"
	"HITS_ToDoList_Tests/internal/domain/enums"
	domainInterfaces "HITS_ToDoList_Tests/internal/domain/interfaces"
	"HITS_ToDoList_Tests/internal/domain/models"
	"HITS_ToDoList_Tests/internal/pkg/utils"
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"strings"
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
	parseTaskName(&name, &deadline, &priority)

	if err := validators.ValidateTask(name, deadline, priority); err != nil {
		return nil, err
	}

	task := models.NewTask(name, description, deadline, nil, priority)

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

	if err := service.tasksRepository.DeleteByID(taskID); err != nil {
		return err
	}

	return nil
}

func (service *TasksServiceImpl) UpdateTask(taskID uuid.UUID, name string, description *string, deadline *time.Time,
	priority *enums.Priority) (*models.Task, error) {
	parseTaskName(&name, &deadline, &priority)
	if err := validators.ValidateTask(name, deadline, priority); err != nil {
		return nil, err
	}

	task, err := service.tasksRepository.GetByID(taskID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, errors.ApplicationError{
			StatusCode: 404,
			Code:       "NotFound",
			Errors:     map[string]string{"message": "Task not found"},
		}
	}

	task.Name = name
	task.Description = description

	if priority != nil {
		task.Priority = *priority
	} else {
		task.Priority = enums.Medium
	}

	task.Deadline = deadline
	if task.Status == enums.Late {
		task.Status = enums.Completed
	} else if task.Status == enums.Overdue {
		task.Status = enums.Active
	}

	task.ChangedAt = utils.Ptr(time.Now())

	if err := service.tasksRepository.Update(*task); err != nil {
		return nil, err
	}

	return task, nil
}

func (service *TasksServiceImpl) ToggleTaskStatus(taskID uuid.UUID, isDone bool) (*models.Task, error) {
	task, err := service.tasksRepository.GetByID(taskID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, errors.ApplicationError{
			StatusCode: 404,
			Code:       "NotFound",
			Errors:     map[string]string{"message": "Task not found"},
		}
	}

	if isDone {
		if task.Deadline != nil && time.Now().After(*task.Deadline) {
			task.Status = enums.Late
		} else {
			task.Status = enums.Completed
		}
	} else {
		if task.Deadline != nil && time.Now().After(*task.Deadline) {
			task.Status = enums.Overdue
		} else {
			task.Status = enums.Active
		}
	}

	task.ChangedAt = utils.Ptr(time.Now())

	if err := service.tasksRepository.Update(*task); err != nil {
		return nil, err
	}

	return task, nil
}

func (service *TasksServiceImpl) UpdateTaskStatuses() {
	tasks, err := service.tasksRepository.GetAll(nil)
	if err != nil {
		fmt.Println("Failed to get all tasks", err.Error())
		return
	}

	curTime := time.Now()
	for _, task := range tasks {
		if task.Status == enums.Active && task.Deadline != nil && curTime.After(*task.Deadline) {
			task.Status = enums.Overdue
			task.ChangedAt = &curTime
			if err := service.tasksRepository.Update(*task); err != nil {
				fmt.Println("Failed to update task", err.Error())
			}
		}
	}
}

func parseTaskName(name *string, deadline **time.Time, priority **enums.Priority) {
	cleanName := *name

	if *deadline == nil {
		pattern := regexp.MustCompile(`!before (\d{2}[.-]\d{2}[.-]\d{4})`)
		matches := pattern.FindStringSubmatch(cleanName)

		if len(matches) == 2 {
			date, err := time.Parse("02.01.2006", strings.ReplaceAll(matches[1], "-", "."))
			if err == nil {
				*deadline = &date
				cleanName = pattern.ReplaceAllString(cleanName, "")
			}
		}
	}

	if *priority == nil {
		patterns := map[string]enums.Priority{
			`!1`: enums.Critical,
			`!2`: enums.High,
			`!3`: enums.Medium,
			`!4`: enums.Low,
		}

		for macro, p := range patterns {
			if strings.Contains(cleanName, macro) {
				*priority = &p
				cleanName = strings.ReplaceAll(cleanName, macro, "")
				break
			}
		}
	}

	*name = strings.TrimSpace(cleanName)
}
