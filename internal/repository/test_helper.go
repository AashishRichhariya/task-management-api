package repository

import (
	"database/sql"
	"os"
	"testing"

	"github.com/AashishRichhariya/task-management-api/internal/database"
)

func SetupTestDB(t *testing.T) *sql.DB {
	// Temporarily override environment for tests
	originalDBName := os.Getenv("DB_NAME")
	os.Setenv("DB_NAME", "taskdb_test")
	
	// Use production connection function with test DB name
	db, err := database.NewPostgresConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	
	// Restore original environment
	if originalDBName != "" {
		os.Setenv("DB_NAME", originalDBName)
	} else {
		os.Unsetenv("DB_NAME")
	}
	
	return db
}

func CleanupTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DELETE FROM tasks")
	if err != nil {
		t.Fatalf("Failed to cleanup test database: %v", err)
	}
}