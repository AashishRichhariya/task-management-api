package repository

import (
	"testing"
	"time"

	"github.com/AashishRichhariya/task-management-api/internal/models"
)

func TestPostgresTaskRepository_CreateTask(t *testing.T) {
	// Setup
	db := SetupTestDB(t)
	defer db.Close()
	defer CleanupTestDB(t, db)
	
	repo := NewPostgresTaskRepository(db)
	
	// Test data
	task := &models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      models.StatusPending,
	}
	
	// Execute
	err := repo.CreateTask(task)
	
	// Assert
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}
	
	// Verify task was created with ID
	if task.ID == 0 {
		t.Error("Expected task ID to be set after creation")
	}
	
	// Verify timestamps were set
	if task.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	
	if task.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestPostgresTaskRepository_GetTaskByID(t *testing.T) {
	// Setup
	db := SetupTestDB(t)
	defer db.Close()
	defer CleanupTestDB(t, db)
	
	repo := NewPostgresTaskRepository(db)
	
	// Create a task first
	originalTask := &models.Task{
		Title:       "Get Test Task",
		Description: "Get Test Description",
		Status:      models.StatusInProgress,
	}
	
	err := repo.CreateTask(originalTask)
	if err != nil {
		t.Fatalf("Failed to create task for test: %v", err)
	}
	
	// Execute
	retrievedTask, err := repo.GetTaskByID(originalTask.ID)
	
	// Assert
	if err != nil {
		t.Fatalf("GetTaskByID failed: %v", err)
	}
	
	if retrievedTask == nil {
		t.Fatal("Expected task to be found")
	}
	
	// Verify all fields match
	if retrievedTask.ID != originalTask.ID {
		t.Errorf("Expected ID %d, got %d", originalTask.ID, retrievedTask.ID)
	}
	
	if retrievedTask.Title != originalTask.Title {
		t.Errorf("Expected Title %s, got %s", originalTask.Title, retrievedTask.Title)
	}
	
	if retrievedTask.Status != originalTask.Status {
		t.Errorf("Expected Status %s, got %s", originalTask.Status, retrievedTask.Status)
	}
}

func TestPostgresTaskRepository_GetTaskByID_NotFound(t *testing.T) {
	// Setup
	db := SetupTestDB(t)
	defer db.Close()
	defer CleanupTestDB(t, db)
	
	repo := NewPostgresTaskRepository(db)
	
	// Execute - try to get non-existent task
	task, err := repo.GetTaskByID(99999)
	
	// Assert
	if err != nil {
		t.Fatalf("GetTaskByID should not return error for missing task: %v", err)
	}
	
	if task != nil {
		t.Error("Expected nil task for non-existent ID")
	}
}

func TestPostgresTaskRepository_GetAllTasks(t *testing.T) {
	// Setup
	db := SetupTestDB(t)
	defer db.Close()
	defer CleanupTestDB(t, db)
	
	repo := NewPostgresTaskRepository(db)
	
	// Create multiple tasks
	tasks := []*models.Task{
		{Title: "Task 1", Status: models.StatusPending},
		{Title: "Task 2", Status: models.StatusCompleted},
		{Title: "Task 3", Status: models.StatusInProgress},
	}
	
	for _, task := range tasks {
		err := repo.CreateTask(task)
		if err != nil {
			t.Fatalf("Failed to create test task: %v", err)
		}
	}
	
	// Execute
	allTasks, totalCount, err := repo.GetAllTasks(10, 1, "", "created_at", "desc")
	
	// Assert
	if totalCount != 3 {
		t.Errorf("Expected total count 3, got %d", totalCount)
	}

	if err != nil {
		t.Fatalf("GetAllTasks failed: %v", err)
	}
	
	if len(allTasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(allTasks))
	}
	
	// Verify ordering (newest first)
	if len(allTasks) >= 2 {
		if allTasks[0].CreatedAt.Before(allTasks[1].CreatedAt) {
			t.Error("Expected tasks to be ordered by created_at DESC")
		}
	}
}

func TestPostgresTaskRepository_UpdateTask(t *testing.T) {
	// Setup
	db := SetupTestDB(t)
	defer db.Close()
	defer CleanupTestDB(t, db)
	
	repo := NewPostgresTaskRepository(db)
	
	// Create task
	task := &models.Task{
		Title:       "Original Title",
		Description: "Original Description",
		Status:      models.StatusPending,
	}
	
	err := repo.CreateTask(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
	
	originalUpdatedAt := task.UpdatedAt
	time.Sleep(time.Millisecond * 10) // Ensure time difference
	
	// Update task
	task.Title = "Updated Title"
	task.Status = models.StatusCompleted
	
	err = repo.UpdateTask(task)
	
	// Assert
	if err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}
	
	// Verify UpdatedAt was changed
	if !task.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
	
	// Verify changes persisted
	updated, err := repo.GetTaskByID(task.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated task: %v", err)
	}
	
	if updated.Title != "Updated Title" {
		t.Errorf("Expected updated title, got %s", updated.Title)
	}
	
	if updated.Status != models.StatusCompleted {
		t.Errorf("Expected updated status, got %s", updated.Status)
	}
}

func TestPostgresTaskRepository_DeleteTask(t *testing.T) {
	// Setup
	db := SetupTestDB(t)
	defer db.Close()
	defer CleanupTestDB(t, db)
	
	repo := NewPostgresTaskRepository(db)
	
	// Create task
	task := &models.Task{
		Title:  "Task to Delete",
		Status: models.StatusPending,
	}
	
	err := repo.CreateTask(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
	
	// Execute delete
	err = repo.DeleteTask(task.ID)
	
	// Assert
	if err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}
	
	// Verify task is gone
	deleted, err := repo.GetTaskByID(task.ID)
	if err != nil {
		t.Fatalf("Error checking if task was deleted: %v", err)
	}
	
	if deleted != nil {
		t.Error("Expected task to be deleted")
	}
}