package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HiroLiang/goat-server/internal/database/repository"
	"github.com/HiroLiang/goat-server/internal/security"
	"github.com/gin-gonic/gin"
)

var userRepo *repository.UserRepo

func InitUserHandler(repo *repository.UserRepo) {
	userRepo = repo
}

func RegisterUserRoutes(r *gin.RouterGroup) {
	{
		r.POST("/login", login)
		r.POST("/register", register)
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

type RegisterRq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// register 處理 POST /api/user/register
func register(c *gin.Context) {
	var rq RegisterRq
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	ip := c.ClientIP()

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// 檢查 是否已有 APPLYING
	exists, err := userRepo.ExistsApplyingByIP(ctx, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal error (check IP)",
		})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "user from this IP is already applying",
		})
		return
	}

	hashedPwd, err := security.HashArgon2Base64(rq.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	//建立 APPLYING user
	u, err := userRepo.CreateApplyingUser(ctx, rq.Email, hashedPwd, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	newName := fmt.Sprintf("user_%d", u.ID)
	if err := userRepo.UpdateUserName(ctx, u.ID, newName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update username",
		})
		return
	}
	u.Name = newName

	//不將 password 傳回去
	c.JSON(http.StatusCreated, gin.H{
		"id":     u.ID,
		"name":   u.Name,
		"email":  u.Email,
		"status": u.UserStatus,
		"ip":     u.UserIP,
	})
}
