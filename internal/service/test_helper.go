package service

import (
	"time"

	"github.com/AashishRichhariya/task-management-api/internal/models"
	"github.com/AashishRichhariya/task-management-api/internal/repository"
)

// Mock repository implementation
type mockTaskRepository struct {
	tasks  map[int]*models.Task
	nextID int
}

func newMockTaskRepository() repository.TaskRepository {
	return &mockTaskRepository{
		tasks:  make(map[int]*models.Task),
		nextID: 1,
	}
}

func (m *mockTaskRepository) CreateTask(task *models.Task) error {
	// Copy the task to avoid pointer issues
	taskCopy := *task
	taskCopy.ID = m.nextID
	m.nextID++
	taskCopy.CreatedAt = time.Now()
	taskCopy.UpdatedAt = time.Now()
	
	// Store the copy
	m.tasks[taskCopy.ID] = &taskCopy
	
	// Update the original task with ID and timestamps
	task.ID = taskCopy.ID
	task.CreatedAt = taskCopy.CreatedAt
	task.UpdatedAt = taskCopy.UpdatedAt
	
	return nil
}

func (m *mockTaskRepository) GetTaskByID(id int) (*models.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, nil // Task not found
	}
	
	// Return a copy to avoid pointer issues
	taskCopy := *task
	return &taskCopy, nil
}

func (m *mockTaskRepository) GetAllTasks(limit, page int, status, sortBy, sortOrder string) ([]models.Task, int, error) {
	// Apply filtering first
	var filteredTasks []*models.Task
	for _, task := range m.tasks {
		if status == "" || string(task.Status) == status {
			filteredTasks = append(filteredTasks, task)
		}
	}
	
	// Convert to slice for return
	var tasks []models.Task
	for _, task := range filteredTasks {
		tasks = append(tasks, *task)
	}
	
	totalCount := len(tasks)
	
	// Apply pagination (simple implementation for testing)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	
	start := (page - 1) * limit
	end := start + limit
	
	if start >= len(tasks) {
		return []models.Task{}, totalCount, nil
	}
	if end > len(tasks) {
		end = len(tasks)
	}
	
	return tasks[start:end], totalCount, nil
}

func (m *mockTaskRepository) UpdateTask(task *models.Task) error {
	_, exists := m.tasks[task.ID]
	if !exists {
		return nil // Simulate sql.ErrNoRows behavior
	}
	
	// Update timestamp
	task.UpdatedAt = time.Now()
	
	// Store updated task
	taskCopy := *task
	m.tasks[task.ID] = &taskCopy
	
	return nil
}

func (m *mockTaskRepository) DeleteTask(id int) error {
	_, exists := m.tasks[id]
	if !exists {
		return nil // Simulate sql.ErrNoRows behavior
	}
	
	delete(m.tasks, id)
	return nil
}
