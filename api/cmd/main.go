package main

import (
	_ "HITS_ToDoList_Tests/docs"
	"HITS_ToDoList_Tests/internal/application/services"
	"HITS_ToDoList_Tests/internal/delivery/handlers"
	"HITS_ToDoList_Tests/internal/delivery/middleware"
	"HITS_ToDoList_Tests/internal/delivery/routes"
	"HITS_ToDoList_Tests/internal/infrastructure/db"
	"HITS_ToDoList_Tests/internal/infrastructure/repositories"
	"HITS_ToDoList_Tests/internal/infrastructure/schedulers"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"time"
)

func main() {
	dbConn, err := db.NewPostgresConnection("localhost", "postgres", "123456", "ToDoDb",
		"5432")
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	if err = db.Migrate(dbConn); err != nil {
		log.Fatalf("Failed to migrate db: %v", err)
	}

	tasksRepository := repositories.NewTasksRepository(dbConn)
	tasksService := services.NewTasksService(tasksRepository)

	schedulers.StartTasksDeadlineScheduling(tasksService, time.Second)

	r := gin.Default()

	// Добавляем CORS middleware первым
	r.Use(middleware.Cors())
	r.Use(middleware.ErrorHandler())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	tasksHandler := handlers.NewTasksHandler(tasksService)
	routes.SetupRoutes(r, tasksHandler)

	r.Run(":8080")

	log.Println("Application started")
}
