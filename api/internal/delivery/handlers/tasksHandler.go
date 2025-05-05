package handlers

import (
	appEnums "HITS_ToDoList_Tests/internal/application/enums"
	"HITS_ToDoList_Tests/internal/application/errors"
	"HITS_ToDoList_Tests/internal/application/interfaces"
	"HITS_ToDoList_Tests/internal/delivery/DTOs"
	"github.com/gin-gonic/gin"
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
// @Router /tasks [post]
func (h *TasksHandler) CreateTask(c *gin.Context) {
	var request DTOs.CreateTaskRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(errors.ApplicationError{
			StatusCode: 400,
			Code:       "invalid_request",
			Errors:     map[string]string{"message": err.Error()},
		})
		return
	}

	task, err := h.tasksService.CreateTask(*request.Name, request.Description, request.Deadline, request.Priority)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"task": task})
}

// GetAllTasks
// @Summary Get all tasks
// @Description Get all tasks with optional sorting
// @Tags tasks
// @Accept json
// @Produce json
// @Param sorting query string false "Sorting" Enums(CreateAsc, CreateDesc, PriorityAsc, PriorityDesc, DeadlineAsc, DeadlineDesc)
// @Success 200 {object} []models.Task
// @Failure 400 {object} errors.ApplicationError "Bad request"
// @Failure 500 {object} map[string]string
// @Router /tasks [get]
func (h *TasksHandler) GetAllTasks(c *gin.Context) {
	sortingParam := c.Query("sorting")
	sorting := appEnums.Sorting(sortingParam)

	tasks, err := h.tasksService.GetAllTasks(&sorting)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}
