package agent

import (
	"github.com/HiroLiang/goat-server/internal/logger"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) bool {
	logger.Log.Error(err.Error())
	switch {

	default:
		_ = c.Error(err)
		return false
	}
}
