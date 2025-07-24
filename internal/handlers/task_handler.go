package handlers

import (
	"net/http"
	"strconv"

	"github.com/AashishRichhariya/task-management-api/internal/models"
	"github.com/AashishRichhariya/task-management-api/internal/service"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService service.TaskServiceInterface
}

func NewTaskHandler(taskService service.TaskServiceInterface) TaskHandlerInterface {
	return &TaskHandler{
		taskService: taskService,
	}
}

// POST /tasks
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	if req.Status == "" {
		req.Status = string(models.StatusPending);
	}

	task, err := h.taskService.CreateTask(req.Title, req.Description, req.Status)
	if err != nil {
		switch err.(type) {
		case service.ValidationError:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Validation failed",
				Message: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to create task",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Task created successfully",
		Data:    task,
	})
}

// GET /tasks/:id
func (h *TaskHandler) GetTask(c *gin.Context) {
	// Get ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid task ID",
			Message: "Task ID must be a valid integer",
		})
		return
	}

	task, err := h.taskService.GetTaskByID(id)
	if err != nil {
		switch err.(type) {
		case service.TaskNotFoundError:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Task not found", 
				Message: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to retrieve task",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Task retrieved successfully",
		Data:    task,
	})
}

// GET /tasks?page=1&limit=10&status=completed
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	status := c.Query("status")          
	sortBy := c.Query("sort_by")   
	sortOrder := c.Query("sort_order")   

	// Convert string parameters to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid page parameter",
			Message: "Page must be a positive integer",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid limit parameter",
			Message: "Limit must be a positive integer",
		})
		return
	}

	// Get paginated tasks from service
	response, err := h.taskService.GetAllTasks(page, limit, status, sortBy, sortOrder)
	if err != nil {
		switch err.(type) {
		case service.ValidationError:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Validation failed",
				Message: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to retrieve tasks",
				Message: err.Error(),
			})
		}
		return
	}

	// Return success response with pagination metadata
	c.IndentedJSON(http.StatusOK, SuccessResponse{
			Message: "Tasks retrieved successfully",
			Data: map[string]any{
					"tasks": response.Tasks,
					"pagination": response.Pagination,
			},
	})
}

// PUT /tasks/:id
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid task ID",
			Message: "Task ID must be a valid integer",
		})
		return
	}

	var req UpdateTaskRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	task, err := h.taskService.UpdateTask(id, req.Title, req.Description, req.Status)
	if err != nil {
		switch err.(type) {
		case service.TaskNotFoundError:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Task not found",
				Message: err.Error(),
			})
		case service.ValidationError:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Validation failed",
				Message: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to update task",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Task updated successfully",
		Data:    task,
	})
}

// DELETE /tasks/:id
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid task ID",
			Message: "Task ID must be a valid integer",
		})
		return
	}

	err = h.taskService.DeleteTask(id)
	if err != nil {
		switch err.(type) {
		case service.TaskNotFoundError:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Task not found",
				Message: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to delete task",
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Task deleted successfully",
	})
}