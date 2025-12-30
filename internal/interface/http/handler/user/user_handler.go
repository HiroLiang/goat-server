package user

import (
	"net/http"

	"github.com/HiroLiang/goat-server/internal/application/user"
	"github.com/HiroLiang/goat-server/internal/interface/http/adapter"
	"github.com/HiroLiang/goat-server/internal/interface/http/middleware"
	"github.com/gin-gonic/gin"
)

// UserHandler Rest api about user
type UserHandler struct {
	userUseCase *user.UseCase
}

// NewUserHandler Create a new UserHandler instance with dependencies
func NewUserHandler(userUseCase *user.UseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// RegisterUserRoutes registers user-related API routes
func (h *UserHandler) RegisterUserRoutes(r *gin.RouterGroup) {
	r.POST("/register", h.register)
	r.POST("/login", h.login)

	r.POST("/logout", middleware.RequireAuthMiddleware(), h.logout)
	r.GET("/me", middleware.RequireAuthMiddleware(), h.getCurrentUser)
}

// @Summary User register
// @Description register with email and password
// @Tags User
// @Accept json
// @Produce json
// @Param payload body RegisterRequest true "Register messages"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/user/register [post]
func (h *UserHandler) register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	data := user.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.userUseCase.Register(c.Request.Context(), adapter.BuildInput(c, data)); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{Message: "Register successful"})
}

// @Summary User Login
// @Description
// Authenticate a user using an email and password.
// If the request already contains an Authorization Bearer token,
// it will be forwarded to the authentication service for session continuity.
// @Tags User
// @Accept json
// @Produce json
// @Param payload body LoginRequest true "Login payload"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/user/login [post]
func (h *UserHandler) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	data := user.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.userUseCase.Login(c.Request.Context(), adapter.BuildInput(c, data))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.Header("Authorization", "Bearer "+output.Token)

	c.JSON(http.StatusOK, LoginResponse{Message: "Login successful"})
}

// @Summary User Logout
// @Description Remove the session token from the redis store.
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/user/logout [post]
func (h *UserHandler) logout(c *gin.Context) {
	err := h.userUseCase.Logout(c.Request.Context(), adapter.BuildEmptyInput(c))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Current User info
// @Description query the current user info
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} CurrentUserResponse
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Not Found"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/user/me [get]
func (h *UserHandler) getCurrentUser(c *gin.Context) {
	output, err := h.userUseCase.CurrentUserInfo(c.Request.Context(), adapter.BuildEmptyInput(c))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, CurrentUserResponse{
		Name:     output.Name,
		Email:    output.Email,
		CreateAt: output.CreateAt,
	})
}
