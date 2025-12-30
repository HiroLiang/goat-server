package adapter

import (
	"net/http"

	"github.com/HiroLiang/goat-server/internal/application/shared"
	"github.com/HiroLiang/goat-server/internal/interface/http/response"
	"github.com/gin-gonic/gin"
)

func MustBaseInput(c *gin.Context) *shared.BaseInput {
	metadata, ok := c.Get("context")
	if !ok {
		c.JSON(http.StatusInternalServerError, response.ErrNotFound("metadata"))
	}
	return metadata.(*shared.BaseInput)
}

func BuildInput[T any](c *gin.Context, data T) shared.UseCaseInput[T] {
	return shared.UseCaseInput[T]{
		Base: *MustBaseInput(c),
		Data: data,
	}
}

func BuildEmptyInput(c *gin.Context) shared.UseCaseInput[struct{}] {
	return shared.UseCaseInput[struct{}]{
		Base: *MustBaseInput(c),
	}
}
