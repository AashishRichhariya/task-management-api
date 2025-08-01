package service

import (
	"testing"

	"github.com/AashishRichhariya/task-management-api/internal/models"
)

// Test functions
func TestTaskService_CreateTask(t *testing.T) {
	// Setup
	mockRepo := newMockTaskRepository()
	service := NewTaskService(mockRepo)
	
	// Test valid task creation
	task, err := service.CreateTask("Test Task", "Description", "pending")
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}
	
	if task.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got %s", task.Title)
	}
	
	if task.ID == 0 {
		t.Error("Expected task ID to be set")
	}
}

func TestTaskService_GetTaskByID(t *testing.T) {
	mockRepo := newMockTaskRepository()
	service := NewTaskService(mockRepo)
	
	// Create a task first
	createdTask, _ := service.CreateTask("Test", "Description", "pending")
	
	// Get the task
	retrievedTask, err := service.GetTaskByID(createdTask.ID)
	if err != nil {
		t.Fatalf("GetTaskByID failed: %v", err)
	}
	
	if retrievedTask.Title != "Test" {
		t.Errorf("Expected title 'Test', got %s", retrievedTask.Title)
	}
}

func TestTaskService_GetTaskByID_NotFound(t *testing.T) {
	mockRepo := newMockTaskRepository()
	service := NewTaskService(mockRepo)
	
	// Test non-existent task
	_, err := service.GetTaskByID(999)
	if err == nil {
		t.Error("Expected TaskNotFoundError")
	}
	
	// Check it's a TaskNotFoundError
	if _, ok := err.(models.TaskNotFoundError); !ok {
		t.Errorf("Expected TaskNotFoundError, got %T", err)
	}
}

func TestTaskService_UpdateTask_Success(t *testing.T) {
	mockRepo := newMockTaskRepository()
	service := NewTaskService(mockRepo)
	
	// Create a task first
	task, _ := service.CreateTask("Original", "Description", "pending")
	
	// Update the task
	updatedTask, err := service.UpdateTask(task.ID, "Updated Title", "Updated Description", "in_progress")
	if err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}
	
	if updatedTask.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %s", updatedTask.Title)
	}
	
	if updatedTask.Status != "in_progress" {
		t.Errorf("Expected status 'in_progress', got %s", updatedTask.Status)
	}
}

func TestTaskService_DeleteTask_Success(t *testing.T) {
	mockRepo := newMockTaskRepository()
	service := NewTaskService(mockRepo)
	
	// Create a task first
	task, _ := service.CreateTask("To Delete", "Description", "pending")
	
	// Delete the task
	err := service.DeleteTask(task.ID)
	if err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}
	
	// Verify task is gone
	_, err = service.GetTaskByID(task.ID)
	if err == nil {
		t.Error("Expected TaskNotFoundError after deletion")
	}
}

func TestTaskService_DeleteTask_NotFound(t *testing.T) {
	mockRepo := newMockTaskRepository()
	service := NewTaskService(mockRepo)
	
	// Try to delete non-existent task
	err := service.DeleteTask(999)
	if err == nil {
		t.Error("Expected TaskNotFoundError")
	}
	
	if _, ok := err.(models.TaskNotFoundError); !ok {
		t.Errorf("Expected TaskNotFoundError, got %T", err)
	}
}

func TestTaskService_GetAllTasks(t *testing.T) {
	mockRepo := newMockTaskRepository()
	service := NewTaskService(mockRepo)
	
	// Create multiple tasks
	service.CreateTask("Task 1", "", "pending")
	service.CreateTask("Task 2", "", "completed")
	service.CreateTask("Task 3", "", "in_progress")
	
	// Get all tasks
	response, err := service.GetAllTasks(1, 10, "", "", "")
	if err != nil {
		t.Fatalf("GetAllTasks failed: %v", err)
	}

	if len(response.Tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(response.Tasks))
	}
}