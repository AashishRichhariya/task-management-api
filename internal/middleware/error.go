package middleware

import (
	"log"
	"net/http"

	"github.com/AashishRichhariya/task-management-api/internal/models"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		if err, ok := recovered.(string); ok {
			log.Printf("Panic recovered: %s", err)
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal server error",
			Message: "Something went wrong. Please try again later.",
		})
	})
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			switch e := err.(type) {
			case models.TaskNotFoundError:
				c.JSON(http.StatusNotFound, models.ErrorResponse{
					Error:   "Task not found",
					Message: e.Error(),
				})
			case models.ValidationError:
				c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Error:   "Validation failed",
					Message: e.Error(),
				})
			case models.BusinessError:
				c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Error:   "Business logic error",
					Message: e.Error(),
				})
			default:
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{
					Error:   "Internal server error",
					Message: "Something went wrong",
				})
			}
			return
		}
	}
}