package user

import (
	"errors"
	"net/http"

	"github.com/HiroLiang/goat-server/internal/application/user"
	"github.com/HiroLiang/goat-server/internal/interface/http/response"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) bool {
	switch {
	case errors.Is(err, user.ErrUserNotFound):
		c.JSON(http.StatusForbidden, response.ErrNotFound("user"))
		return true

	case errors.Is(err, user.ErrInvalidUser):
		c.JSON(http.StatusForbidden, response.ErrInvalid("user"))
		return true

	case errors.Is(err, user.ErrInvalidPassword):
		c.JSON(http.StatusForbidden, response.ErrInvalid("password"))
		return true

	case errors.Is(err, user.ErrInvalidEmail):
		c.JSON(http.StatusForbidden, response.ErrInvalid("email"))
		return true

	case errors.Is(err, user.ErrGenerateToken):
		c.JSON(http.StatusForbidden, response.ErrAuthFailed)
		return true

	default:
		return false
	}
}
