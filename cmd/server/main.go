package main

import (
	"log"

	"github.com/AashishRichhariya/task-management-api/internal/database"
	"github.com/AashishRichhariya/task-management-api/internal/handlers"
	"github.com/AashishRichhariya/task-management-api/internal/repository"
	"github.com/AashishRichhariya/task-management-api/internal/service"
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


	log.Println("Starting server on :8080")
	router.Run(":8080")
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
			tasks.POST("", taskHandler.CreateTask)   
			tasks.GET("/:id", taskHandler.GetTask)          
			tasks.GET("", taskHandler.GetAllTasks)        
			tasks.PUT("/:id", taskHandler.UpdateTask)    
			tasks.DELETE("/:id", taskHandler.DeleteTask)  
		}
	}
	
	return router
}