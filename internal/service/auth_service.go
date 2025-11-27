package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/HiroLiang/goat-server/internal/database/model"
	"github.com/HiroLiang/goat-server/internal/database/repository"
	"github.com/HiroLiang/goat-server/internal/database/session"
	"github.com/HiroLiang/goat-server/internal/security"
	"github.com/HiroLiang/goat-server/internal/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotActivated   = errors.New("user account is not activated")
	ErrAlreadyLoggedIn    = errors.New("user is already logged in")
	ErrInvalidToken       = errors.New("invalid token")
	ErrSessionNotFound    = errors.New("session not found or expired")
)

// AuthService 認證服務
type AuthService struct {
	userRepo         *repository.UserRepo
	rateLimitService *RateLimitService
	sessionManager   *session.SessionManager
	jwtSecret        string
	jwtExpiration    int // hours
}

// NewAuthService 創建新的認證服務
func NewAuthService(
	userRepo *repository.UserRepo,
	rateLimitService *RateLimitService,
	sessionManager *session.SessionManager,
	jwtSecret string,
	jwtExpiration int,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		rateLimitService: rateLimitService,
		sessionManager:   sessionManager,
		jwtSecret:        jwtSecret,
		jwtExpiration:    jwtExpiration,
	}
}

// LoginRequest 登入請求
type LoginRequest struct {
	Email    string
	Password string
	IP       string
	Token    string // 可選：現有的 session token
}

// LoginResponse 登入響應
type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// RegisterRequest 註冊請求
type RegisterRequest struct {
	Email    string
	Password string
	IP       string
}

// RegisterResponse 註冊響應
type RegisterResponse struct {
	User *model.User `json:"user"`
}

// Login 用戶登入
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// 使用 defer 統一處理：釋放並發登入鎖和記錄登入嘗試
	loginSuccess := false
	defer func() {
		_ = s.rateLimitService.ReleaseConcurrentLoginLock(ctx, req.Email)
		_ = s.rateLimitService.RecordLoginAttempt(ctx, req.IP, req.Email, loginSuccess)
	}()

	// 0. 檢查是否已有有效的 session token
	if req.Token != "" {
		claims, err := utils.ValidateToken(req.Token, s.jwtSecret)
		if err == nil && claims.SessionID != "" {
			if s.sessionManager.IsSessionValid(ctx, claims.SessionID) {
				return nil, ErrAlreadyLoggedIn
			}
		}
	}

	// 1. 檢查 Rate Limit
	err := s.rateLimitService.CheckLoginAttempt(ctx, req.IP, req.Email)
	if err != nil {
		// 記錄失敗嘗試（因為被限流，需要在 defer 之前處理）
		_ = s.rateLimitService.RecordLoginAttempt(ctx, req.IP, req.Email, false)
		return nil, err
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if s.userRepo.IsDisabled(user) {
		return nil, repository.ErrUserDisabled
	}

	if !s.userRepo.IsActivated(user) {
		return nil, ErrUserNotActivated
	}

	// 驗證密碼（使用 Argon2）
	if !security.VerifyArgon2Base64(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// 創建 session（先創建 session 獲取 session ID）
	sessionData := session.SessionData{
		UserID: user.ID,
		Email:  user.Email,
		UserIP: req.IP,
	}
	sessionID, err := s.sessionManager.CreateSession(ctx, sessionData)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// 生成 JWT token（包含 session ID）
	token, err := utils.GenerateToken(user.ID, user.Email, sessionID, s.jwtSecret, s.jwtExpiration)
	if err != nil {
		// 如果生成 token 失敗，刪除已創建的 session
		_ = s.sessionManager.DeleteSession(ctx, sessionID)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	err = s.userRepo.UpdateUserIP(ctx, user.ID, req.IP)
	if err != nil {
		fmt.Printf("failed to update user IP: %v\n", err)
	}

	loginSuccess = true

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// Register 用戶註冊
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	// 1. 檢查 Rate Limit
	err := s.rateLimitService.CheckRegisterAttempt(ctx, req.IP, req.Email)
	if err != nil {
		// 記錄註冊嘗試
		_ = s.rateLimitService.RecordRegisterAttempt(ctx, req.IP, req.Email)
		return nil, err
	}

	// 2. 檢查 email 是否已存在
	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email exists: %w", err)
	}
	if exists {
		// 記錄註冊嘗試
		_ = s.rateLimitService.RecordRegisterAttempt(ctx, req.IP, req.Email)
		return nil, repository.ErrUserAlreadyExists
	}

	// 3. 加密密碼（使用 Argon2）
	hashedPassword, err := security.HashArgon2Base64(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 4. 創建用戶（user_status = ACTIVATED）
	user := &model.User{
		Name:       req.Email, // 如果沒有提供 name，使用 email 作為默認值
		Email:      req.Email,
		Password:   hashedPassword,
		UserStatus: model.UserStatusActivated,
		UserIP:     req.IP,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		// 記錄註冊嘗試
		_ = s.rateLimitService.RecordRegisterAttempt(ctx, req.IP, req.Email)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 記錄成功的註冊嘗試
	_ = s.rateLimitService.RecordRegisterAttempt(ctx, req.IP, req.Email)

	return &RegisterResponse{
		User: user,
	}, nil
}

// validateTokenAndGetSessionID 驗證 token 並提取 session ID（共用函數）
func (s *AuthService) validateTokenAndGetSessionID(ctx context.Context, token string, checkSession bool) (string, *utils.Claims, error) {
	// 驗證 token
	claims, err := utils.ValidateToken(token, s.jwtSecret)
	if err != nil {
		return "", nil, ErrInvalidToken
	}

	if claims.SessionID == "" {
		return "", nil, fmt.Errorf("session ID not found in token")
	}

	// 如果需要，檢查 session 是否有效
	if checkSession {
		if !s.sessionManager.IsSessionValid(ctx, claims.SessionID) {
			return "", nil, ErrSessionNotFound
		}
	}

	return claims.SessionID, claims, nil
}

// Logout 用戶登出
func (s *AuthService) Logout(ctx context.Context, token string) error {
	sessionID, _, err := s.validateTokenAndGetSessionID(ctx, token, false)
	if err != nil {
		return err
	}

	return s.sessionManager.DeleteSession(ctx, sessionID)
}

// RefreshTokenRequest 刷新 token 請求
type RefreshTokenRequest struct {
	Token string
}

// RefreshTokenResponse 刷新 token 響應
type RefreshTokenResponse struct {
	Token string `json:"token"`
}

// RefreshToken 刷新 token
func (s *AuthService) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*RefreshTokenResponse, error) {
	// 驗證 token 並檢查 session 有效性
	sessionID, claims, err := s.validateTokenAndGetSessionID(ctx, req.Token, true)
	if err != nil {
		return nil, err
	}

	// 刷新 session（獲取新的 session ID）
	newSessionID, err := s.sessionManager.RefreshSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh session: %w", err)
	}

	// 生成新 token（使用新的 session ID）
	newToken, err := utils.GenerateToken(claims.UserID, claims.Email, newSessionID, s.jwtSecret, s.jwtExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &RefreshTokenResponse{
		Token: newToken,
	}, nil
}

// GetCurrentUser 獲取當前用戶信息
func (s *AuthService) GetCurrentUser(ctx context.Context, token string) (*model.User, error) {
	// 驗證 token 並檢查 session 有效性
	sessionID, claims, err := s.validateTokenAndGetSessionID(ctx, token, true)
	if err != nil {
		return nil, err
	}

	// 自動刷新 session 過期時間（用戶活動時延長 session）
	err = s.sessionManager.RefreshSessionExpiration(ctx, sessionID)
	if err != nil {
		// 如果刷新失敗，記錄但不影響獲取用戶信息
		// 因為 session 仍然有效，只是過期時間沒有延長
	}

	// 獲取用戶信息
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
