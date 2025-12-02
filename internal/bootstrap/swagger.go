package bootstrap

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterSwaggerRoutes(swagger *gin.RouterGroup) {
	swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
