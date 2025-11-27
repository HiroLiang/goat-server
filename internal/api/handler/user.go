package handler

import (
	"fmt"
	"context"
	"net/http"
	"strings"

	"github.com/HiroLiang/goat-server/internal/database"
	"github.com/HiroLiang/goat-server/internal/database/repository"
	"github.com/HiroLiang/goat-server/internal/database/session"
	"github.com/HiroLiang/goat-server/internal/service"
	"github.com/HiroLiang/goat-server/internal/security"
	"github.com/gin-gonic/gin"
)

var authService *service.AuthService

// InitUserHandler 初始化 User Handler
func InitUserHandler(jwtSecret string, jwtExpiration int) {
	// 創建依賴
	userRepo := repository.NewUserRepo(database.Postgres)
	rateLimitService := service.NewRateLimitService(database.RedisClient, nil)
	sessionManager := session.NewSessionManager(database.RedisClient, "session", jwtExpiration)

	// 創建認證服務
	authService = service.NewAuthService(
		userRepo,
		rateLimitService,
		sessionManager,
		jwtSecret,
		jwtExpiration,
	)
}

var userRepo *repository.UserRepo

func RegisterUserRoutes(r *gin.RouterGroup) {
	r.POST("/login", login)
	r.POST("/register", register)
	r.POST("/logout", logout)
	r.POST("/refresh-token", refreshToken)
	r.GET("/me", getCurrentUser)
}

// LoginRequest 登入請求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	{
		r.POST("/login", login)
		r.POST("/register", register)
	}
}

// LoginResponse 登入響應
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// @Summary 用戶登入
// @Description 用戶使用 email 和密碼登入
// @Tags 用戶
// @Accept json
// @Produce json
// @Param body body LoginRequest true "登入信息"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 429 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/user/login [post]
func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// 獲取客戶端 IP
	clientIP := c.ClientIP()

	// 檢查 Authorization header 中是否已有 token
	var existingToken string
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			existingToken = parts[1]
		}
	}

	// 調用認證服務
	resp, err := authService.Login(c.Request.Context(), service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
		IP:       clientIP,
		Token:    existingToken,
	})

	if err != nil {
		handleLoginError(c, err)
		return
	}

	// 設置 Authorization header
	c.Header("Authorization", "Bearer "+resp.Token)

	c.JSON(http.StatusOK, gin.H{
		"token": "Bearer " + resp.Token,
	})
}

// RegisterRequest 註冊請求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterResponse 註冊響應
type RegisterResponse struct {
	Message string      `json:"message"`
	User    interface{} `json:"user"`
}

// @Summary 用戶註冊
// @Description 用戶使用 email 和密碼註冊新帳號
// @Tags 用戶
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "註冊信息"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 429 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/user/register [post]
func register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// 獲取客戶端 IP
	clientIP := c.ClientIP()

	// 調用認證服務
	resp, err := authService.Register(c.Request.Context(), service.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		IP:       clientIP,
	})

	if err != nil {
		handleRegisterError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"user":    resp.User,
	})
}

// handleLoginError 處理登入錯誤
func handleLoginError(c *gin.Context, err error) {
	switch err {
	case service.ErrInvalidCredentials:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
	case service.ErrUserNotActivated:
		c.JSON(http.StatusForbidden, gin.H{"error": "User account is not activated"})
	case repository.ErrUserDisabled:
		c.JSON(http.StatusForbidden, gin.H{"error": "User account is disabled"})
	case service.ErrAlreadyLoggedIn:
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already logged in"})
	case service.ErrRateLimitExceeded:
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many login attempts, please try again later"})
	case service.ErrConcurrentLogin:
		c.JSON(http.StatusConflict, gin.H{"error": "Concurrent login detected, please try again"})
	case service.ErrUserLocked:
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "User account is locked due to too many failed attempts"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": err.Error()})
	}
}

// handleAuthError 處理認證相關錯誤（logout, refreshToken, getCurrentUser）
func handleAuthError(c *gin.Context, err error) {
	switch err {
	case service.ErrInvalidToken:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
	case service.ErrSessionNotFound:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired or invalid"})
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed", "details": err.Error()})
	}
}

// handleRegisterError 處理註冊錯誤
func handleRegisterError(c *gin.Context, err error) {
	switch err {
	case repository.ErrUserAlreadyExists:
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
	case service.ErrRateLimitExceeded:
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many registration attempts, please try again later"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": err.Error()})
	}
}

// getTokenFromHeader 從 Authorization header 中提取 token
func getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is required")
	}

	// 移除前後空格
	authHeader = strings.TrimSpace(authHeader)

	// 檢查是否以 Bearer 開頭（大小寫不敏感）
	authHeaderLower := strings.ToLower(authHeader)
	if !strings.HasPrefix(authHeaderLower, "bearer ") {
		// 如果沒有 Bearer 前綴，嘗試直接使用整個 header 作為 token（向後兼容）
		// 但先檢查是否包含空格，如果包含空格則可能是格式錯誤
		if strings.Contains(authHeader, " ") {
			return "", fmt.Errorf("invalid authorization format: expected 'Bearer <token>' but got '%s'", authHeader)
		}
		// 如果沒有空格，假設整個 header 就是 token
		return authHeader, nil
	}

	// 提取 Bearer 後面的 token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid authorization format: missing token after 'Bearer'")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", fmt.Errorf("invalid authorization format: token is empty")
	}

	return token, nil
}

// @Summary 用戶登出
// @Description 用戶登出，刪除當前 session
// @Tags 用戶
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/user/logout [post]
func logout(c *gin.Context) {
	token, err := getTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = authService.Logout(c.Request.Context(), token)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// RefreshTokenResponse 刷新 token 響應
type RefreshTokenResponse struct {
	Token string `json:"token"`
}

// @Summary 刷新 Token
// @Description 刷新用戶的 JWT token
// @Tags 用戶
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} RefreshTokenResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/user/refresh-token [post]
func refreshToken(c *gin.Context) {
	// 從 Authorization header 獲取 token
	token, err := getTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp, err := authService.RefreshToken(c.Request.Context(), service.RefreshTokenRequest{
		Token: token,
	})

	if err != nil {
		handleAuthError(c, err)
		return
	}

	// 設置新的 Authorization header
	c.Header("Authorization", "Bearer "+resp.Token)

	c.JSON(http.StatusOK, gin.H{
		"token": "Bearer " + resp.Token,
	})
}

// @Summary 獲取當前用戶信息
// @Description 獲取當前登入用戶的信息
// @Tags 用戶
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/user/me [get]
func getCurrentUser(c *gin.Context) {
	token, err := getTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := authService.GetCurrentUser(c.Request.Context(), token)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
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
