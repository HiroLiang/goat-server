package handler

import (
	"net/http"

	"github.com/HiroLiang/goat-server/internal/logger"
	"github.com/gin-gonic/gin"
)

func RegisterTestRoutes(r *gin.RouterGroup) {
	r.GET("", test)
}

// @Summary Test api is ready
// @Router /api/test [get]
func test(context *gin.Context) {
	logger.Log.Info("Test OK")
	context.JSON(http.StatusOK, gin.H{"result": "Test OK"})
}
