package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AashishRichhariya/task-management-api/internal/middleware"
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

func (m *MockTaskService) GetAllTasks(page, limit int, status, sortBy, sortOrder string) (*models.PaginatedTasksResponse, error) {
	args := m.Called(page, limit, status, sortBy, sortOrder)
	if response := args.Get(0); response != nil {
		return response.(*models.PaginatedTasksResponse), args.Error(1)
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
	
	requestBody := models.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
	}
	jsonBody, _ := json.Marshal(requestBody)
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/tasks", append(middleware.ValidateCreateTaskBody(), handler.CreateTask)...)
	
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	
	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response models.SuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task created successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestCreateTask_ValidationError(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)
	
	requestBody := models.CreateTaskRequest{
		Title:  "Test Task",
		Status: "invalid",
	}
	jsonBody, _ := json.Marshal(requestBody)
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/tasks", append(middleware.ValidateCreateTaskBody(), handler.CreateTask)...)
	
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetTask_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)
	
	task := &models.Task{ID: 1, Title: "Test Task", Status: models.StatusPending}
	mockService.On("GetTaskByID", 1).Return(task, nil)
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/:id", append(middleware.ValidateTaskID(), handler.GetTask)...)
	
	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response models.SuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task retrieved successfully", response.Message)
	mockService.AssertExpectations(t)
}

func TestGetTask_NotFound(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)
	
	mockService.On("GetTaskByID", 999).Return(nil, service.TaskNotFoundError{ID: 999})
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/:id", append(middleware.ValidateTaskID(), handler.GetTask)...)
	
	req, _ := http.NewRequest("GET", "/tasks/999", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	
	assert.Equal(t, http.StatusNotFound, recorder.Code)

	// Verify error response structure
	var errorResponse models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Task not found", errorResponse.Error)
	mockService.AssertExpectations(t)
}

func TestGetTask_InvalidID(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks/:id", append(middleware.ValidateTaskID(), handler.GetTask)...)
	
	req, _ := http.NewRequest("GET", "/tasks/abc", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	
	var errorResponse models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Contains(t, errorResponse.Error, "Invalid ID parameter")
}

func TestGetAllTasks_Success(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)
	
	tasks := []models.Task{
		{ID: 1, Title: "Task 1", Status: models.StatusPending},
		{ID: 2, Title: "Task 2", Status: models.StatusCompleted},
	}
	
	paginatedResponse := &models.PaginatedTasksResponse{
		Tasks: tasks,
		Pagination: models.PaginationMeta{
			Page:    1,
			Limit:   10,
			Total:   2,
			Pages:   1,
			HasNext: false,
			HasPrev: false,
		},
	}
	
	mockService.On("GetAllTasks", 1, 10, "", "created_at", "desc").Return(paginatedResponse, nil)
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/tasks", append(middleware.ValidateTaskQuery(), handler.GetAllTasks)...)
	
	req, _ := http.NewRequest("GET", "/tasks", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	
	assert.Equal(t, http.StatusOK, recorder.Code)
	
	var response models.SuccessResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Tasks retrieved successfully", response.Message)
	
	mockService.AssertExpectations(t)
}