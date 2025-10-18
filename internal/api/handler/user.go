package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	{
		r.POST("login", login)
	}
}

func login(c *gin.Context) {
	var rq LoginRq
	if err := c.ShouldBind(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}

type LoginRq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
