package routes

import (
	"github.com/HiroLiang/goat-server/internal/api/handler"
	"github.com/HiroLiang/goat-server/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRestRoutes(r *gin.Engine) {
	group := r.Group("/api", middleware.ErrorHandler())
	{
		handler.RegisterTestRoutes(group.Group("/test"))
		handler.RegisterUserRoutes(group.Group("/user"))
	}
}
