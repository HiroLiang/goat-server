package device

import (
	"net/http"

	"github.com/HiroLiang/goat-server/internal/interface/http/adapter"
	"github.com/HiroLiang/goat-server/internal/logger"
	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
}

func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{}
}

func (h *DeviceHandler) RegisterDeviceRoutes(r *gin.RouterGroup) {
	r.POST("/register", h.registerDeviceId)
}

// @Summary registerDeviceId
// @Description try to register device id checking is id already registered
// @Tags Device
// @Accept json
// @Produce json
// @Param payload body RegisterDeviceIdRequest true "Register device"
// @Success 200 {object} RegisterDeviceIdResponse
// @Router /api/device/register [post]
func (h *DeviceHandler) registerDeviceId(c *gin.Context) {
	var req RegisterDeviceIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	input := adapter.BuildInput(c, req)

	logger.Log.Info("Register device id: " + req.DeviceID)

	if input.Data.DeviceID != "" {
		c.JSON(http.StatusOK, RegisterDeviceIdResponse{
			Success:  true,
			DeviceID: input.Data.DeviceID,
			Message:  "Device registered successfully",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
		})
	}
}
