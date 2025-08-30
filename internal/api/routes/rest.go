package routes

import (
	"github.com/HiroLiang/goat-chat-server/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRestRoutes(r *gin.Engine) {
	group := r.Group("/api")
	{
		handler.RegisterRoutes(group.Group("/test"))
	}
}
