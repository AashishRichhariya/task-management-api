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

// GET /tasks
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.taskService.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve tasks",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Tasks retrieved successfully",
		Data:    tasks,
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