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

		// Processing request
		context.Next()

		// Handling error
		if len(context.Errors) > 0 {
			err := context.Errors[0].Err

			// Check if the error is ApiError (known error)
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

			// Unknown error, return 500
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    http.StatusInternalServerError,
					"message": err.Error(),
				},
			})
		}
	}
}
