package service

import (
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

func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	return s.taskRepo.GetAllTasks()
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
