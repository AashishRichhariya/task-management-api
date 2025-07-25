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

type TaskServiceInterface interface {
	CreateTask(title, description, status string) (*models.Task, error)
	GetTaskByID(id int) (*models.Task, error)
	GetAllTasks(page, limit int, status, sortBy, sortOrder string) (*models.PaginatedTasksResponse, error)
	UpdateTask(id int, title, description, status string) (*models.Task, error)
	DeleteTask(id int) error
}


func NewTaskService(taskRepo repository.TaskRepository) TaskServiceInterface {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

func (s *TaskService) CreateTask(title, description, status string) (*models.Task, error) {
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
		return nil, models.TaskNotFoundError{ID: id}
	}
	
	return task, nil
}

func (s *TaskService) GetAllTasks(page, limit int, status, sortBy, sortOrder string) (*models.PaginatedTasksResponse, error) {
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


	pagination := models.PaginationMeta{
		Page:    page,
		Limit:   limit,
		Total:   totalCount,
		Pages:   totalPages,
		HasNext: page < totalPages,
		HasPrev: page > 1,
	}

	return &models.PaginatedTasksResponse{
		Tasks:      tasks,
		Pagination: pagination,
	}, nil
}

func (s *TaskService) UpdateTask(id int, title, description, status string) (*models.Task, error) {
	existingTask, err := s.GetTaskByID(id)
	if err != nil {
		return nil, err
	}

	if title != "" {
		existingTask.Title = strings.TrimSpace(title)
	}
	if description != "" {
		existingTask.Description = strings.TrimSpace(description)  
	}
	if status != "" {
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
