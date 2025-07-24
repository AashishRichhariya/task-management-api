package service

import (
	"fmt"

	"github.com/AashishRichhariya/task-management-api/internal/models"
)

type TaskServiceInterface interface {
	CreateTask(title, description, status string) (*models.Task, error)
	GetTaskByID(id int) (*models.Task, error)
	GetAllTasks(page, limit int, status, sortBy, sortOrder string) (*PaginatedTasksResponse, error)
	UpdateTask(id int, title, description, status string) (*models.Task, error)
	DeleteTask(id int) error
}

type PaginationMeta struct {
	Page      int  `json:"page"`
	Limit     int  `json:"limit"`
	Total     int  `json:"total"`
	Pages     int  `json:"pages"`
	HasNext   bool `json:"has_next"`
	HasPrev   bool `json:"has_prev"`
}

type PaginatedTasksResponse struct {
	Tasks      []models.Task   `json:"tasks"`
	Pagination PaginationMeta `json:"pagination"`
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