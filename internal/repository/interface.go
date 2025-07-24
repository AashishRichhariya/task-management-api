package repository

import "github.com/AashishRichhariya/task-management-api/internal/models"

type TaskRepository interface {
	// Create operations
	CreateTask(task *models.Task) error
	
	// Read operations  
	GetTaskByID(id int) (*models.Task, error)
	GetAllTasks() ([]models.Task, error)
	
	// Update operations
	UpdateTask(task *models.Task) error
	
	// Delete operations
	DeleteTask(id int) error
}