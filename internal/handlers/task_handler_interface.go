package handlers

import (
	"github.com/gin-gonic/gin"
)

type TaskHandlerInterface interface {
	CreateTask(c *gin.Context)
	GetTask(c *gin.Context)
	GetAllTasks(c *gin.Context)
	UpdateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    any `json:"data,omitempty"`
}
