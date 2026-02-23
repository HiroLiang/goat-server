package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) RegisterHealthRoues(r *gin.RouterGroup) {
	r.GET("", h.healthCheck)
}

// @Summary api health check
// @Description api health check
// @Tags Health
// @Router /api/health [get]
func (h *HealthHandler) healthCheck(c *gin.Context) {
	c.Status(http.StatusOK)
}
