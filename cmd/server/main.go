package main

import (
	"log"

	"github.com/AashishRichhariya/task-management-api/internal/database"
	"github.com/AashishRichhariya/task-management-api/internal/handlers"
	"github.com/AashishRichhariya/task-management-api/internal/middleware"
	"github.com/AashishRichhariya/task-management-api/internal/repository"
	"github.com/AashishRichhariya/task-management-api/internal/service"
	"github.com/AashishRichhariya/task-management-api/internal/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// Database setup
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	
	// Dependency injection
	taskRepo := repository.NewPostgresTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)
	
	// Router setup
	router := setupRoutes(taskHandler)

	port := utils.GetEnv("APP_PORT", "8080")
	log.Println("Starting server on :" + port)
	router.Run(":" + port)
}

func setupRoutes(taskHandler handlers.TaskHandlerInterface) *gin.Engine {
	router := gin.Default()
	
	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)
	
	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Task routes
		tasks := v1.Group("/tasks")
		{
			tasks.POST("", append(middleware.ValidateCreateTaskBody(), taskHandler.CreateTask)...)   
			tasks.GET("/:id", append(middleware.ValidateTaskID(), taskHandler.GetTask)...)          
			tasks.GET("", append(middleware.ValidateTaskQuery(), taskHandler.GetAllTasks)...)        
			tasks.PUT("/:id", append(
				append(middleware.ValidateTaskID(), middleware.ValidateUpdateTaskBody()...), 
					taskHandler.UpdateTask,
				)...)
			tasks.DELETE("/:id", append(middleware.ValidateTaskID(), taskHandler.DeleteTask)...)  
		}
	}	
	return router
}