package middleware

import (
	"errors"
	"net/http"
	appErrors "sportlink/api/app/errors"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			var appErr appErrors.AppError
			if errors.As(err.Err, &appErr) {
				status := http.StatusInternalServerError
				switch appErr.Code {
				case appErrors.InvalidRequestDataErrorCode:
					status = http.StatusBadRequest
				case appErrors.NotFoundErrorCode:
					status = http.StatusNotFound
				case appErrors.UnauthorizedErrorCode:
					status = http.StatusUnauthorized
				}

				c.AbortWithStatusJSON(status, gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    appErrors.UnexpectedErrorCode,
				"message": "Oops, something went wrong",
			})
		}
	}
}
