package bootstrap

import (
	"github.com/HiroLiang/goat-server/internal/config"
	"github.com/HiroLiang/goat-server/internal/interface/http/handler/test"
	"github.com/HiroLiang/goat-server/internal/interface/http/handler/user"
	"github.com/HiroLiang/goat-server/internal/interface/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRestRoutes(group *gin.RouterGroup, useCases *UseCases, dependencies *Dependencies) {

	// Global middleware
	group.Use(middleware.ErrorHandler())
	group.Use(middleware.AuthMiddleware(dependencies.TokenService))
	group.Use(middleware.ContextMiddleware())

	// Test Handler
	if config.Env("APP_ENV", "dev") == "dev" {
		var testHandler = test.NewTestHandler()
		testHandler.RegisterTestRoutes(group.Group("/test"))
	}

	// User Handler
	var userHandler = user.NewUserHandler(useCases.UserUseCase)
	userHandler.RegisterUserRoutes(group.Group("/user"))

}
