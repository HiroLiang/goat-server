package middleware

import (
	"net/http"
	"strings"

	"github.com/HiroLiang/goat-server/internal/database/session"
	"github.com/HiroLiang/goat-server/internal/utils"
	"github.com/gin-gonic/gin"
)

// JWTMiddleware JWT 認證中間件（包含 Session 檢查）
func JWTMiddleware(jwtSecret string, sessionManager *session.SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 header 獲取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 檢查 Bearer 前綴
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := parts[1]

		// 驗證 token
		claims, err := utils.ValidateToken(token, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 檢查 session 是否有效（如果 token 包含 session ID）
		if sessionManager != nil && claims.SessionID != "" {
			if !sessionManager.IsSessionValid(c.Request.Context(), claims.SessionID) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired or invalid"})
				c.Abort()
				return
			}

			// 自動刷新 session 過期時間（用戶活動時延長 session）
			_ = sessionManager.RefreshSessionExpiration(c.Request.Context(), claims.SessionID)
		}

		// 將用戶信息存入 context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("token", token)
		if claims.SessionID != "" {
			c.Set("session_id", claims.SessionID)
		}

		c.Next()
	}
}
