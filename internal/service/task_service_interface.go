package service

import (
	"fmt"

	"github.com/AashishRichhariya/task-management-api/internal/models"
)

type TaskServiceInterface interface {
	CreateTask(title, description, status string) (*models.Task, error)
	GetTaskByID(id int) (*models.Task, error)
	GetAllTasks() ([]models.Task, error)
	UpdateTask(id int, title, description, status string) (*models.Task, error)
	DeleteTask(id int) error
}

type TaskNotFoundError struct {
	ID int
}

func (e TaskNotFoundError) Error() string {
	return fmt.Sprintf("task with id %d not found", e.ID)
}

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}