package middleware

import (
	"net/http"
	"strings"

	"github.com/AashishRichhariya/task-management-api/internal/models"
	"github.com/gin-gonic/gin"
)

func ValidateTaskID() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(c *gin.Context) {
			var param models.TaskIDParam
			
			if err := c.ShouldBindUri(&param); err != nil {
				c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{
					Error:   "Invalid ID parameter",
					Message: err.Error(),
				})
				c.Abort()
				return
			}
			
			// Store validated ID in context
			c.Set("taskID", param.ID)
			c.Next()
		},
	}
}

func ValidateTaskQuery() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(c *gin.Context) {
			var query models.TaskQueryParams
			
			if err := c.ShouldBindQuery(&query); err != nil {
				c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{
					Error:   "Invalid query parameters",
					Message: err.Error(),
				})
				c.Abort()
				return
			}
			
			query.SetDefaults()
			
			// Store in context
			c.Set("taskQuery", query)
			c.Next()
		},
	}
}

func ValidateCreateTaskBody() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(c *gin.Context) {
			var req models.CreateTaskRequest
			
			if err := c.ShouldBindJSON(&req); err != nil {
				c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{
					Error:   "Invalid request body",
					Message: err.Error(),
				})
				c.Abort()
				return
			}
			
			// Set default status and trim whitespace
			if strings.TrimSpace(req.Status) == "" {
				req.Status = string(models.StatusPending)
			}
			req.Title = strings.TrimSpace(req.Title)
			req.Description = strings.TrimSpace(req.Description)
			
			// Store in context
			c.Set("createTaskReq", req)
			c.Next()
		},
	}
}

func ValidateUpdateTaskBody() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(c *gin.Context) {
			var req models.UpdateTaskRequest
			
			if err := c.ShouldBindJSON(&req); err != nil {
				c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{
					Error:   "Invalid request body",
					Message: err.Error(),
				})
				c.Abort()
				return
			}
			
			// Trim whitespace for non-empty fields
			if req.Title != "" {
				req.Title = strings.TrimSpace(req.Title)
			}
			if req.Description != "" {
				req.Description = strings.TrimSpace(req.Description)
			}
			
			// Store in context
			c.Set("updateTaskReq", req)
			c.Next()
		},
	}
}

// Helper functions for handlers to extract validated data
func GetTaskID(c *gin.Context) int {
	return c.MustGet("taskID").(int)
}

func GetTaskQuery(c *gin.Context) models.TaskQueryParams {
	return c.MustGet("taskQuery").(models.TaskQueryParams)
}

func GetCreateTaskRequest(c *gin.Context) models.CreateTaskRequest {
	return c.MustGet("createTaskReq").(models.CreateTaskRequest)
}

func GetUpdateTaskRequest(c *gin.Context) models.UpdateTaskRequest {
	return c.MustGet("updateTaskReq").(models.UpdateTaskRequest)
}