package test

import (
	"net/http"

	"github.com/HiroLiang/goat-server/internal/logger"
	"github.com/gin-gonic/gin"
)

type TestHandler struct{}

func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

func (h *TestHandler) RegisterTestRoutes(r *gin.RouterGroup) {
	r.GET("", h.test)
}

// @Summary Test api is ready
// @Router /api/test [get]
func (h *TestHandler) test(context *gin.Context) {
	logger.Log.Info("Test OK")
	context.JSON(http.StatusOK, gin.H{"result": "Test OK"})
}
