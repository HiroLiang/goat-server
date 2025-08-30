package handler

import (
	"net/http"

	"github.com/HiroLiang/goat-chat-server/internal/logger"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/", test)
}

func test(context *gin.Context) {
	logger.Log.Info("Test OK")
	context.JSON(http.StatusOK, gin.H{"result": "Test OK"})
}
