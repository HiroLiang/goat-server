package middleware

import (
	"github.com/HiroLiang/goat-server/internal/logger"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		logger.Log.Info(token)
	}
}
