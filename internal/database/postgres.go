package database

import (
	"database/sql"
	"fmt"

	"github.com/AashishRichhariya/task-management-api/internal/utils"
	_ "github.com/lib/pq" // PostgreSQL driver registration
)

func NewPostgresConnection() (*sql.DB, error) {
	// Get connection details from environment variables
	host := utils.GetEnv("DB_HOST", "localhost")
	port := utils.GetEnv("DB_PORT", "5432")
	user := utils.GetEnv("DB_USER", "postgres")
	password := utils.GetEnv("DB_PASSWORD", "password")
	dbname := utils.GetEnv("DB_NAME", "taskdb")
	
	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	
	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return db, nil
}