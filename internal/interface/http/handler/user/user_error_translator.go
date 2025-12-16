package user

import (
	"errors"
	"net/http"

	"github.com/HiroLiang/goat-server/internal/domain/user"
	"github.com/HiroLiang/goat-server/internal/interface/http/response"
	"github.com/HiroLiang/goat-server/internal/logger"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	logger.Log.Error(err.Error())
	switch {
	case errors.Is(err, user.ErrUserNotFound):
		c.JSON(http.StatusNotFound, response.ErrNotFound("user"))
		return

	case errors.Is(err, user.ErrInvalidUser):
		c.JSON(http.StatusForbidden, response.ErrInvalid("user"))
		return

	case errors.Is(err, user.ErrInvalidPassword):
		c.JSON(http.StatusUnauthorized, response.ErrInvalid("password"))
		return

	case errors.Is(err, user.ErrInvalidEmail):
		c.JSON(http.StatusBadRequest, response.ErrInvalid("email"))
		return

	case errors.Is(err, user.ErrGenerateToken):
		c.JSON(http.StatusInternalServerError, response.ErrAuthFailed)
		return

	default:
		_ = c.Error(err)
		return
	}
}
