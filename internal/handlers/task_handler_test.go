package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AashishRichhariya/task-management-api/internal/models"
	"github.com/AashishRichhariya/task-management-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Simple mock for testing handlers
type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(title, description, status string) (*models.Task, error) {
	args := m.Called(title, description, status)
	if task := args.Get(0); task != nil {
		return task.(*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskService) GetTaskByID(id int) (*models.Task, error) {
	args := m.Called(id)
	if task := args.Get(0); task != nil {
		return task.(*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskService) GetAllTasks(page, limit int, status, sortBy, sortOrder string) (*service.PaginatedTasksResponse, error) {
	args := m.Called(page, limit, status, sortBy, sortOrder)
	if response := args.Get(0); response != nil {
		return response.(*service.PaginatedTasksResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskService) UpdateTask(id int, title, description, status string) (*models.Task, error) {
	args := m.Called(id, title, description, status)
	if task := args.Get(0); task != nil {
		return task.(*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskService) DeleteTask(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateTask_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)
	
	task := &models.Task{
		ID:          1,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      models.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	mockService.On("CreateTask", "Test Task", "Test Description", "pending").Return(task, nil)
	
	requestBody := CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
	}
	jsonBody, _ := json.Marshal(requestBody)
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/tasks", handler.CreateTask)
	
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	
	assert.Equal(t, http.StatusCreated, recorder.Code)
	mockService.AssertExpectations(t)
}

func TestCreateTask_ValidationError(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)
	
	mockService.On("CreateTask", "Test Task", "", "invalid").Return(nil, service.ValidationError{
		Field: "status", Message: "invalid status",
	})
	
	requestBody := CreateTaskRequest{
		Title:  "Test Task",
		Status: "invalid",
	}
	jsonBody, _ := json.Marshal(requestBody)
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/tasks", handler.CreateTask)
	
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	mockService.AssertExpectations(t)
}

func TestGetTask_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)
	
	task := &models.Task{ID: 1, Title: "Test Task", Status: models.StatusPending}
	mockService.On("GetTaskByID", 1).Return(task, nil)
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/:id", handler.GetTask)
	
	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	
	assert.Equal(t, http.StatusOK, recorder.Code)
	mockService.AssertExpectations(t)
}

func TestGetTask_NotFound(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)
	
	mockService.On("GetTaskByID", 999).Return(nil, service.TaskNotFoundError{ID: 999})
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/:id", handler.GetTask)
	
	req, _ := http.NewRequest("GET", "/tasks/999", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	
	assert.Equal(t, http.StatusNotFound, recorder.Code)
	mockService.AssertExpectations(t)
}