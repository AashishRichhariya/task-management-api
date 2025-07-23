package main

import (
	"log"

	"github.com/AashishRichhariya/task-management-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	
	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)
	router.GET("/tasks", handlers.GetTasks)


	log.Println("Starting server on :8080")
	router.Run(":8080")
}