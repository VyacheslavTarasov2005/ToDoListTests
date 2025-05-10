package handlers

import (
	appEnums "HITS_ToDoList_Tests/internal/application/enums"
	"HITS_ToDoList_Tests/internal/application/errors"
	"HITS_ToDoList_Tests/internal/application/interfaces"
	"HITS_ToDoList_Tests/internal/delivery/DTOs"
	"HITS_ToDoList_Tests/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// @BasePath /tasks

type TasksHandler struct {
	tasksService interfaces.TasksService
}

func NewTasksHandler(tasksService interfaces.TasksService) *TasksHandler {
	return &TasksHandler{tasksService: tasksService}
}

// CreateTask
// @Summary Create a task
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body DTOs.CreateTaskRequest true "Task"
// @Success 201 {object} models.Task
// @Failure 400 {object} errors.ApplicationError "Bad request"
// @Failure 500 "Internal server error"
// @Router /tasks [post]
func (h *TasksHandler) CreateTask(c *gin.Context) {
	var request DTOs.CreateTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(errors.ApplicationError{
			StatusCode: 400,
			Code:       "InvalidRequest",
			Errors:     map[string]string{"message": err.Error()},
		})
		return
	}

	task, err := h.tasksService.CreateTask(*request.Name, request.Description, request.Deadline, request.Priority)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, DTOs.TaskResponse{
		ID:          task.ID,
		CreatedAt:   task.CreatedAt,
		ChangedAt:   task.ChangedAt,
		Name:        task.Name,
		Description: task.Description,
		Deadline:    task.Deadline,
		Status:      task.Status,
		Priority:    task.Priority,
	})
}

// GetAllTasks
// @Summary Get all tasks
// @Description Get all tasks with optional sorting
// @Tags tasks
// @Accept json
// @Produce json
// @Param sorting query string false "Sorting" Enums(CreateAsc, CreateDesc, PriorityAsc, PriorityDesc, DeadlineAsc, DeadlineDesc)
// @Success 200 {object} []models.Task
// @Failure 500 "Internal server error"
// @Router /tasks [get]
func (h *TasksHandler) GetAllTasks(c *gin.Context) {
	var sorting = utils.Ptr(c.Query("sorting"))
	if *sorting == "" {
		sorting = nil
	}

	if sorting != nil {
		err := appEnums.ValidateSorting(appEnums.Sorting(*sorting))
		if err != nil {
			c.Error(errors.ApplicationError{
				StatusCode: 400,
				Code:       "InvalidRequest",
				Errors:     map[string]string{"message": err.Error()},
			})
			return
		}
	}

	tasks, err := h.tasksService.GetAllTasks((*appEnums.Sorting)(sorting))
	if err != nil {
		c.Error(err)
		return
	}

	response := make([]DTOs.TaskResponse, len(tasks))
	for i, item := range tasks {
		response[i] = DTOs.TaskResponse{
			ID:          item.ID,
			CreatedAt:   item.CreatedAt,
			ChangedAt:   item.ChangedAt,
			Name:        item.Name,
			Description: item.Description,
			Deadline:    item.Deadline,
			Status:      item.Status,
			Priority:    item.Priority,
		}
	}

	c.JSON(http.StatusOK, response)
}

// DeleteTask
// @Summary Delete task
// @Description Delete task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 204 "No Content"
// @Failure 400 {object} errors.ApplicationError "Bad request"
// @Failure 404 {object} errors.ApplicationError "Not found"
// @Failure 500 "Internal server error"
// @Router /tasks/{id} [delete]
func (h *TasksHandler) DeleteTask(c *gin.Context) {
	taskIDParam := c.Param("id")

	taskID, err := uuid.Parse(taskIDParam)
	if err != nil {
		c.Error(errors.ApplicationError{
			StatusCode: 400,
			Code:       "InvalidRequest",
			Errors:     map[string]string{"message": err.Error()},
		})
		return
	}

	err = h.tasksService.DeleteTask(taskID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateTask
// @Summary Update task
// @Description Update task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param task body DTOs.UpdateTaskRequest true "Task"
// @Success 200 {object} DTOs.TaskResponse
// @Failure 400 {object} errors.ApplicationError "Bad request"
// @Failure 404 {object} errors.ApplicationError "Not found"
// @Failure 500 "Internal server error"
// @Router /tasks/{id} [put]
func (h *TasksHandler) UpdateTask(c *gin.Context) {
	taskIDParam := c.Param("id")

	taskID, err := uuid.Parse(taskIDParam)
	if err != nil {
		c.Error(errors.ApplicationError{
			StatusCode: 400,
			Code:       "InvalidRequest",
			Errors:     map[string]string{"message": err.Error()},
		})
		return
	}

	var request DTOs.UpdateTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(errors.ApplicationError{
			StatusCode: 400,
			Code:       "InvalidRequest",
			Errors:     map[string]string{"message": err.Error()},
		})
		return
	}

	task, err := h.tasksService.UpdateTask(taskID, *request.Name, request.Description, request.Deadline,
		request.Priority)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, DTOs.TaskResponse{
		ID:          task.ID,
		CreatedAt:   task.CreatedAt,
		ChangedAt:   task.ChangedAt,
		Name:        task.Name,
		Description: task.Description,
		Deadline:    task.Deadline,
		Status:      task.Status,
		Priority:    task.Priority,
	})
}

// ToggleTaskStatus
// @Summary Toggle task's status
// @Description Change task's status
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param task body DTOs.ToggleTaskStatusRequest true "Task"
// @Success 200 {object} DTOs.TaskResponse
// @Failure 400 {object} errors.ApplicationError "Bad request"
// @Failure 404 {object} errors.ApplicationError "Not found"
// @Failure 500 "Internal server error"
// @Router /tasks/{id}/toggle [patch]
func (h *TasksHandler) ToggleTaskStatus(c *gin.Context) {
	taskIDParam := c.Param("id")

	taskID, err := uuid.Parse(taskIDParam)
	if err != nil {
		c.Error(errors.ApplicationError{
			StatusCode: 400,
			Code:       "InvalidRequest",
			Errors:     map[string]string{"message": err.Error()},
		})
		return
	}

	var request DTOs.ToggleTaskStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(errors.ApplicationError{
			StatusCode: 400,
			Code:       "InvalidRequest",
			Errors:     map[string]string{"message": err.Error()},
		})
		return
	}

	task, err := h.tasksService.ToggleTaskStatus(taskID, *request.IsDone)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, DTOs.TaskResponse{
		ID:          task.ID,
		CreatedAt:   task.CreatedAt,
		ChangedAt:   task.ChangedAt,
		Name:        task.Name,
		Description: task.Description,
		Deadline:    task.Deadline,
		Status:      task.Status,
		Priority:    task.Priority,
	})
}
