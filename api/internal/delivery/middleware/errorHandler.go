package middleware

import (
	"HITS_ToDoList_Tests/internal/application/errors"
	defaultErrors "errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		var appErr errors.ApplicationError
		if defaultErrors.As(c.Errors[0].Err, &appErr) {
			c.AbortWithStatusJSON(appErr.StatusCode, gin.H{
				"code":   appErr.Code,
				"errors": appErr.Errors,
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
	}
}
