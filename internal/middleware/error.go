package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Code    int
	Message string
}

func (e *ApiError) Error() string {
	return e.Message
}

func ErrorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()

		if len(context.Errors) > 0 {
			err := context.Errors[0].Err

			var apiErr *ApiError
			if errors.As(err, &apiErr) {
				context.JSON(apiErr.Code, gin.H{
					"error": gin.H{
						"code":    apiErr.Code,
						"message": apiErr.Message,
					},
				})
				return
			}

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    http.StatusInternalServerError,
					"message": err.Error(),
				},
			})
		}
	}
}
