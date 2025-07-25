package handlers

import (
	"net/http"

	"github.com/AashishRichhariya/task-management-api/internal/middleware"
	"github.com/AashishRichhariya/task-management-api/internal/models"
	"github.com/AashishRichhariya/task-management-api/internal/service"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService service.TaskServiceInterface
}

type TaskHandlerInterface interface {
	CreateTask(c *gin.Context)
	GetTask(c *gin.Context)
	GetAllTasks(c *gin.Context)
	UpdateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
}


func NewTaskHandler(taskService service.TaskServiceInterface) TaskHandlerInterface {
	return &TaskHandler{
		taskService: taskService,
	}
}

// POST /tasks
func (h *TaskHandler) CreateTask(c *gin.Context) {
	req := middleware.GetCreateTaskRequest(c)

	task, err := h.taskService.CreateTask(req.Title, req.Description, req.Status)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create task",
			Message: err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, models.SuccessResponse{
		Message: "Task created successfully",
		Data:    task,
	})
}

// GET /tasks/:id
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := middleware.GetTaskID(c)

	task, err := h.taskService.GetTaskByID(id)
	if err != nil {
		switch err.(type) {
		case service.TaskNotFoundError:
			c.IndentedJSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Task not found", 
				Message: err.Error(),
			})
		default:
			c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Failed to retrieve task",
				Message: err.Error(),
			})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, models.TaskResponse{
		Message: "Task retrieved successfully",
		Data:    task,
	})
}

// GET /tasks?page=1&limit=10&status=completed
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	query := middleware.GetTaskQuery(c)
	
	response, err := h.taskService.GetAllTasks(
		query.Page,
		query.Limit,
		query.Status,
		query.SortBy,
		query.SortOrder,
	)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to retrieve tasks",
			Message: err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, models.TasksResponse{
		Message: "Tasks retrieved successfully",
		Data:    response,
	})
}

// PUT /tasks/:id
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := middleware.GetTaskID(c)
	req := middleware.GetUpdateTaskRequest(c)
	
	task, err := h.taskService.UpdateTask(id, req.Title, req.Description, req.Status)
	if err != nil {
		switch err.(type) {
		case service.TaskNotFoundError:
			c.IndentedJSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Task not found",
				Message: err.Error(),
			})
		default:
			c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Failed to update task",
				Message: err.Error(),
			})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, models.SuccessResponse{
		Message: "Task updated successfully",
		Data:    task,
	})
}

// DELETE /tasks/:id
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := middleware.GetTaskID(c)
	
	err := h.taskService.DeleteTask(id)
	if err != nil {
		switch err.(type) {
		case service.TaskNotFoundError:
			c.IndentedJSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Task not found",
				Message: err.Error(),
			})
		default:
			c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Failed to delete task",
				Message: err.Error(),
			})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, models.SuccessResponse{
		Message: "Task deleted successfully",
	})
}