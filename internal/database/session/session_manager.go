package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

// SessionData 存儲在 Redis 中的會話數據
type SessionData struct {
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	UserIP    string    `json:"user_ip"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SessionManager 會話管理器
type SessionManager struct {
	client        *redis.Client
	sessionPrefix string
	expiration    time.Duration
}

// NewSessionManager 創建新的會話管理器
func NewSessionManager(client *redis.Client, sessionPrefix string, expirationHours int) *SessionManager {
	return &SessionManager{
		client:        client,
		sessionPrefix: sessionPrefix,
		expiration:    time.Hour * time.Duration(expirationHours),
	}
}

// GenerateSessionID 生成唯一的 session ID
func (sm *SessionManager) GenerateSessionID() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate session ID: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CreateSession 創建新的會話，返回 session ID
func (sm *SessionManager) CreateSession(ctx context.Context, data SessionData) (string, error) {
	// 生成唯一的 session ID
	sessionID, err := sm.GenerateSessionID()
	if err != nil {
		return "", err
	}

	data.CreatedAt = time.Now()
	data.ExpiresAt = time.Now().Add(sm.expiration)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal session data: %w", err)
	}

	key := sm.getSessionKey(sessionID)
	err = sm.client.Set(ctx, key, jsonData, sm.expiration).Err()
	if err != nil {
		return "", fmt.Errorf("failed to save session: %w", err)
	}

	return sessionID, nil
}

// GetSession 獲取會話數據
func (sm *SessionManager) GetSession(ctx context.Context, sessionID string) (*SessionData, error) {
	key := sm.getSessionKey(sessionID)

	result, err := sm.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrSessionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var data SessionData
	err = json.Unmarshal([]byte(result), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}

	// 檢查是否過期
	if time.Now().After(data.ExpiresAt) {
		sm.DeleteSession(ctx, sessionID) // 清除過期會話
		return nil, ErrSessionExpired
	}

	return &data, nil
}

// DeleteSession 刪除會話
func (sm *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	key := sm.getSessionKey(sessionID)
	err := sm.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// RefreshSession 刷新會話過期時間，返回新的 session ID
func (sm *SessionManager) RefreshSession(ctx context.Context, sessionID string) (string, error) {
	data, err := sm.GetSession(ctx, sessionID)
	if err != nil {
		return "", err
	}

	// 刪除舊 session
	_ = sm.DeleteSession(ctx, sessionID)

	// 更新過期時間並創建新 session
	data.ExpiresAt = time.Now().Add(sm.expiration)
	return sm.CreateSession(ctx, *data)
}

// RefreshSessionExpiration 僅刷新會話過期時間（不創建新 session ID，性能更好）
func (sm *SessionManager) RefreshSessionExpiration(ctx context.Context, sessionID string) error {
	key := sm.getSessionKey(sessionID)

	// 檢查 session 是否存在
	exists, err := sm.client.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check session: %w", err)
	}
	if exists == 0 {
		return ErrSessionNotFound
	}

	// 延長過期時間
	err = sm.client.Expire(ctx, key, sm.expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to refresh session expiration: %w", err)
	}

	return nil
}

// IsSessionValid 檢查會話是否存在且有效（使用 Redis EXISTS 優化性能）
func (sm *SessionManager) IsSessionValid(ctx context.Context, sessionID string) bool {
	key := sm.getSessionKey(sessionID)
	exists, err := sm.client.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	return exists > 0
}

// getSessionKey 獲取 Redis 鍵名
func (sm *SessionManager) getSessionKey(sessionID string) string {
	return fmt.Sprintf("%s:%s", sm.sessionPrefix, sessionID)
}

// GetUserSessions 獲取用戶的所有會話（用於管理多設備登入）
// 注意：此方法需要額外的索引結構，暫時留空
func (sm *SessionManager) GetUserSessions(ctx context.Context, userID int64) ([]string, error) {
	// TODO: 實現用戶會話索引
	return nil, errors.New("not implemented yet")
}

// DeleteUserSessions 刪除用戶的所有會話（用於登出所有設備）
// 注意：此方法需要額外的索引結構，暫時留空
func (sm *SessionManager) DeleteUserSessions(ctx context.Context, userID int64) error {
	// TODO: 實現批量刪除用戶會話
	return errors.New("not implemented yet")
}
