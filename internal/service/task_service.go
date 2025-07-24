package service

import (
	"fmt"
	"strings"

	"github.com/AashishRichhariya/task-management-api/internal/models"
	"github.com/AashishRichhariya/task-management-api/internal/repository"
)

type TaskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskServiceInterface {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

func (s *TaskService) CreateTask(title, description, status string) (*models.Task, error) {
	// Business validation
	if err := s.validateTitle(title); err != nil {
		return nil, err
	}
	
	if err := s.validateDescription(description); err != nil {
		return nil, err
	}
	
	if err := s.validateStatus(status); err != nil {
		return nil, err
	}
	
	// Set default status if empty
	if status == "" {
		status = string(models.StatusPending)
	}
	
	// Create task model
	task := &models.Task{
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Status:      models.TaskStatus(status),
	}
	
	// Delegate to repository
	err := s.taskRepo.CreateTask(task)
	if err != nil {
		return nil, err
	}
	
	return task, nil
}

func (s *TaskService) GetTaskByID(id int) (*models.Task, error) {
	task, err := s.taskRepo.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	
	if task == nil {
		return nil, TaskNotFoundError{ID: id}
	}
	
	return task, nil
}

func (s *TaskService) GetAllTasks(page, limit int, status, sortBy, sortOrder string) (*PaginatedTasksResponse, error) {
	// Validate and set defaults
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	
	// Validate status if provided
	if status != "" {
		taskStatus := models.TaskStatus(status)
		if !taskStatus.IsValid() {
			return nil, ValidationError{
				Field:   "status",
				Message: "invalid status value",
			}
		}
	}
	
	validSortFields := map[string]bool{
		"id":         true,
		"title":      true,
		"status":     true,
		"created_at": true,
		"updated_at": true,
	}
	
	if sortBy != "" && !validSortFields[sortBy] {
		return nil, ValidationError{
			Field:   "sort_by",
			Message: "invalid sort field",
		}
	}

	// Get tasks from repository
	tasks, totalCount, err := s.taskRepo.GetAllTasks(limit, page, status, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	// Calculate pagination metadata
	totalPages := (totalCount + limit - 1) / limit 
	if totalPages == 0 {
		totalPages = 1
	}

	pagination := PaginationMeta{
		Page:    page,
		Limit:   limit,
		Total:   totalCount,
		Pages:   totalPages,
		HasNext: page < totalPages,
		HasPrev: page > 1,
	}

	return &PaginatedTasksResponse{
		Tasks:      tasks,
		Pagination: pagination,
	}, nil
}

func (s *TaskService) UpdateTask(id int, title, description, status string) (*models.Task, error) {
	existingTask, err := s.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	
	// Validate new values if provided
	if title != "" {
		if err := s.validateTitle(title); err != nil {
			return nil, err
		}
		existingTask.Title = strings.TrimSpace(title)
	}
	
	if description != "" {
		if err := s.validateDescription(description); err != nil {
			return nil, err
		}
		existingTask.Description = strings.TrimSpace(description)
	}
	
	if status != "" {
		if err := s.validateStatus(status); err != nil {
			return nil, err
		}
		existingTask.Status = models.TaskStatus(status)
	}
	
	// Update in repository
	err = s.taskRepo.UpdateTask(existingTask)
	if err != nil {
		return nil, err
	}
	
	return existingTask, nil
}

func (s *TaskService) DeleteTask(id int) error {
	// Check if task exists
	_, err := s.GetTaskByID(id)
	if err != nil {
		return err 
	}
	
	// Delete from repository
	return s.taskRepo.DeleteTask(id)
}
