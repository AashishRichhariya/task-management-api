package repository

import "github.com/AashishRichhariya/task-management-api/internal/models"

type TaskRepository interface {
	// Create operations
	CreateTask(task *models.Task) error
	
	// Read operations  
	GetTaskByID(id int) (*models.Task, error)
	GetAllTasks(limit, page int, status, sortBy, sortOrder string) ([]models.Task, int, error)
	
	// Update operations
	UpdateTask(task *models.Task) error
	
	// Delete operations
	DeleteTask(id int) error
}