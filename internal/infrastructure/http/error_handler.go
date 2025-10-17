package http // O donde sea que organices tus middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/javito2003/shortener_url/internal/app"
)

// ErrorHandler is a Gin middleware that handles errors and sends appropriate HTTP responses.

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		lastError := c.Errors.Last().Err

		var appErr *app.AppError
		if errors.As(lastError, &appErr) {
			switch appErr.Type {
			case app.ErrNotFound:
				c.JSON(http.StatusNotFound, gin.H{"message": appErr.Message})
			case app.ErrUnauthorized:
				c.JSON(http.StatusUnauthorized, gin.H{"message": appErr.Message})
			case app.ErrConflict:
				c.JSON(http.StatusConflict, gin.H{"message": appErr.Message})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "an unexpected application error occurred"})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "an internal server error occurred"})
		}
	}
}
