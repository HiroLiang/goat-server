package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/HiroLiang/goat-server/internal/application/shared"
	"github.com/HiroLiang/goat-server/internal/application/shared/auth"
	"github.com/HiroLiang/goat-server/internal/interface/http/response"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware try to validate auth token from the header
func AuthMiddleware(tokenService auth.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get auth token from the header
		authHeader := c.GetHeader("Authorization")

		// Validate token if exists
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			if session, err := tokenService.Validate(c.Request.Context(), token); err == nil {
				c.Set("authContext", &shared.AuthContext{
					UserID: session.UserID,
					Token:  token,
				})
			}
		}

		c.Next()

		if strings.HasPrefix(authHeader, "Bearer ") {
			c.Header("Authorization", authHeader)
		}
	}
}

// RequireAuthMiddleware require auth context to be set or abort with 401
func RequireAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Check if auth context exists
		if _, exists := c.Get("authContext"); !exists {
			fmt.Println("authContext not exists")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": response.ErrAuthFailed,
			})
			return
		}
		c.Next()
	}
}
