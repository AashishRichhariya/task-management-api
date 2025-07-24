package service

import (
	"strings"

	"github.com/AashishRichhariya/task-management-api/internal/models"
)

func (s *TaskService) validateTitle(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return ValidationError{
			Field:   "title",
			Message: "title is required",
		}
	}
	
	if len(title) > 255 {
		return ValidationError{
			Field:   "title",
			Message: "title must be less than 255 characters",
		}
	}
	
	return nil
}

func (s *TaskService) validateDescription(description string) error {
	if len(description) > 1000 {
		return ValidationError{
			Field:   "description",
			Message: "description must be less than 1000 characters",
		}
	}
	return nil
}

func (s *TaskService) validateStatus(status string) error {
	if status == "" {
		return nil // Empty status is allowed (defaults to pending)
	}
	
	taskStatus := models.TaskStatus(status)
	if !taskStatus.IsValid() {
		return ValidationError{
			Field:   "status",
			Message: "status must be one of: pending, in_progress, completed, closed",
		}
	}
	
	return nil
}