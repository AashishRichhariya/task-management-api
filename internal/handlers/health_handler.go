package handlers

import (
	"net/http"
	"os"

	"github.com/AashishRichhariya/task-management-api/internal/database"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Check database connection
	dbStatus := "connected"
	db, err := database.NewPostgresConnection()
	if err != nil {
		dbStatus = "disconnected"
	} else {
		if err := db.Ping(); err != nil {
			dbStatus = "disconnected"
		}
		db.Close()
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status":   "healthy",
		"service":  "task-management-api",
		"instance": hostname,
		"database": dbStatus,
	})
}