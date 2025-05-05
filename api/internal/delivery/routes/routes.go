package routes

import (
	"HITS_ToDoList_Tests/internal/delivery/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, tasksHandler *handlers.TasksHandler) {
	tasks := router.Group("/tasks")
	{
		tasks.POST("", tasksHandler.CreateTask)
		tasks.GET("", tasksHandler.GetAllTasks)
		tasks.DELETE("/:id", tasksHandler.DeleteTask)
		tasks.PATCH("/:id", tasksHandler.UpdateTask)
	}
}
