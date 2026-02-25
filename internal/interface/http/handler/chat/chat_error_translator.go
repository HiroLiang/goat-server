package chat

import (
	"errors"
	"net/http"

	"github.com/HiroLiang/goat-server/internal/domain/chatgroup"
	"github.com/HiroLiang/goat-server/internal/domain/user"
	"github.com/HiroLiang/goat-server/internal/interface/http/response"
	"github.com/HiroLiang/goat-server/internal/logger"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	logger.Log.Error(err.Error())
	switch {
	case errors.Is(err, chatgroup.ErrNotFound):
		c.JSON(http.StatusNotFound, response.ErrNotFound("chat group"))
		return

	case errors.Is(err, chatgroup.ErrDeleted):
		c.JSON(http.StatusGone, response.ErrorResponse{
			Code:    "CHAT_GROUP_DELETED",
			Message: "this chat group no longer exists",
		})
		return

	case errors.Is(err, chatgroup.ErrForbidden):
		c.JSON(http.StatusForbidden, response.ErrInvalid("chat group access"))
		return

	case errors.Is(err, user.ErrInvalidUser):
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Code:    "INVALID_USER",
			Message: "invalid user identity",
		})
		return

	default:
		_ = c.Error(err)
		return
	}
}
