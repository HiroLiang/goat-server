package middleware

import (
	"github.com/HiroLiang/goat-server/internal/application/shared"
	"github.com/gin-gonic/gin"
)

// ContextMiddleware Build context data
func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get auth context
		var authCtx *shared.AuthContext
		if v, ok := c.Get("authContext"); ok {
			authCtx = v.(*shared.AuthContext)
		}

		// Set context
		c.Set("context", &shared.BaseInput{
			Request: shared.RequestContext{
				IP:      c.ClientIP(),
				TraceID: c.GetHeader("traceparent"),
			},
			Auth: authCtx,
		})

		c.Next()
	}
}
