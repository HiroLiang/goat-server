package middleware

import (
	"github.com/HiroLiang/goat-server/internal/application/auth"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO Process
		c.Next()
		// TODO Process
	}
}
