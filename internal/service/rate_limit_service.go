package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrUserLocked        = errors.New("user is locked due to too many login attempts")
	ErrConcurrentLogin   = errors.New("concurrent login detected")
)

// RateLimitConfig Rate Limit 配置
type RateLimitConfig struct {
	// 全局限制：60次/分鐘
	GlobalLimit  int
	GlobalWindow time.Duration

	// IP 限制：5次/分鐘
	IPLimit  int
	IPWindow time.Duration

	// 用戶限制：5次/分鐘
	UserLimit  int
	UserWindow time.Duration

	// 用戶鎖定時間
	UserLockDuration time.Duration
}

// DefaultRateLimitConfig 默認配置
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		GlobalLimit:      60,
		GlobalWindow:     time.Minute,
		IPLimit:          5,
		IPWindow:         time.Minute,
		UserLimit:        5,
		UserWindow:       time.Minute,
		UserLockDuration: 15 * time.Minute,
	}
}

// RateLimitService Rate Limit 服務
type RateLimitService struct {
	client *redis.Client
	config *RateLimitConfig
}

// NewRateLimitService 創建新的 Rate Limit 服務
func NewRateLimitService(client *redis.Client, config *RateLimitConfig) *RateLimitService {
	if config == nil {
		config = DefaultRateLimitConfig()
	}
	return &RateLimitService{
		client: client,
		config: config,
	}
}

func (s *RateLimitService) CheckLoginAttempt(ctx context.Context, ip string, email string) error {
	if locked, err := s.isUserLocked(ctx, email); err != nil {
		return err
	} else if locked {
		return ErrUserLocked
	}

	if err := s.checkConcurrentLogin(ctx, email); err != nil {
		return err
	}

	if err := s.checkGlobalLimit(ctx); err != nil {
		return err
	}

	if err := s.checkIPLimit(ctx, ip); err != nil {
		return err
	}

	if err := s.checkUserLimit(ctx, email); err != nil {
		return err
	}

	return nil
}

// RecordLoginAttempt 記錄登入嘗試
func (s *RateLimitService) RecordLoginAttempt(ctx context.Context, ip string, email string, success bool) error {
	// 記錄全局計數
	if err := s.incrementCounter(ctx, "login:global", s.config.GlobalWindow); err != nil {
		return err
	}

	// 記錄 IP 計數
	if err := s.incrementCounter(ctx, fmt.Sprintf("login:ip:%s", ip), s.config.IPWindow); err != nil {
		return err
	}

	// 記錄用戶計數
	if err := s.incrementCounter(ctx, fmt.Sprintf("login:user:%s", email), s.config.UserWindow); err != nil {
		return err
	}

	// 如果登入失敗，記錄失敗次數
	if !success {
		if err := s.recordFailedAttempt(ctx, email); err != nil {
			return err
		}
	} else {
		// 登入成功，清除失敗記錄
		if err := s.clearFailedAttempts(ctx, email); err != nil {
			return err
		}
	}

	return nil
}

// ReleaseConcurrentLoginLock 釋放並發登入鎖
func (s *RateLimitService) ReleaseConcurrentLoginLock(ctx context.Context, email string) error {
	key := fmt.Sprintf("login:lock:%s", email)
	return s.client.Del(ctx, key).Err()
}

// checkGlobalLimit 檢查全局限制
func (s *RateLimitService) checkGlobalLimit(ctx context.Context) error {
	count, err := s.getCounter(ctx, "login:global")
	if err != nil {
		return err
	}

	if count >= s.config.GlobalLimit {
		return ErrRateLimitExceeded
	}

	return nil
}

// checkIPLimit 檢查 IP 限制
func (s *RateLimitService) checkIPLimit(ctx context.Context, ip string) error {
	key := fmt.Sprintf("login:ip:%s", ip)
	count, err := s.getCounter(ctx, key)
	if err != nil {
		return err
	}

	if count >= s.config.IPLimit {
		return ErrRateLimitExceeded
	}

	return nil
}

// checkUserLimit 檢查用戶限制
func (s *RateLimitService) checkUserLimit(ctx context.Context, email string) error {
	key := fmt.Sprintf("login:user:%s", email)
	count, err := s.getCounter(ctx, key)
	if err != nil {
		return err
	}

	if count >= s.config.UserLimit {
		return ErrRateLimitExceeded
	}

	return nil
}

func (s *RateLimitService) checkConcurrentLogin(ctx context.Context, email string) error {
	key := fmt.Sprintf("login:lock:%s", email)

	success, err := s.client.SetNX(ctx, key, "1", 30*time.Second).Result()
	if err != nil {
		return err
	}

	if !success {
		return ErrConcurrentLogin
	}

	return nil
}

func (s *RateLimitService) isUserLocked(ctx context.Context, email string) (bool, error) {
	key := fmt.Sprintf("login:locked:%s", email)
	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (s *RateLimitService) recordFailedAttempt(ctx context.Context, email string) error {
	key := fmt.Sprintf("login:failed:%s", email)

	// 增加失敗計數
	count, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return err
	}

	// 設置過期時間
	s.client.Expire(ctx, key, s.config.UserWindow)

	// 如果失敗次數達到用戶限制，鎖定用戶
	if count >= int64(s.config.UserLimit) {
		lockKey := fmt.Sprintf("login:locked:%s", email)
		err = s.client.Set(ctx, lockKey, "1", s.config.UserLockDuration).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

// clearFailedAttempts 清除失敗嘗試記錄
func (s *RateLimitService) clearFailedAttempts(ctx context.Context, email string) error {
	key := fmt.Sprintf("login:failed:%s", email)
	return s.client.Del(ctx, key).Err()
}

// incrementCounter 增加計數器
func (s *RateLimitService) incrementCounter(ctx context.Context, key string, window time.Duration) error {
	// 增加計數
	count, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return err
	}

	// 如果是第一次計數，設置過期時間
	if count == 1 {
		s.client.Expire(ctx, key, window)
	}

	return nil
}

// getCounter 獲取計數器值
func (s *RateLimitService) getCounter(ctx context.Context, key string) (int, error) {
	result, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	var count int
	_, err = fmt.Sscanf(result, "%d", &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetRemainingAttempts 獲取剩餘嘗試次數（用於提示用戶）
func (s *RateLimitService) GetRemainingAttempts(ctx context.Context, email string) (int, error) {
	key := fmt.Sprintf("login:user:%s", email)
	count, err := s.getCounter(ctx, key)
	if err != nil {
		return 0, err
	}

	remaining := s.config.UserLimit - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// CheckRegisterAttempt 檢查註冊嘗試的限流
func (s *RateLimitService) CheckRegisterAttempt(ctx context.Context, ip string, email string) error {
	// 檢查全局限制（重用登入的全局限制）
	if err := s.checkGlobalLimit(ctx); err != nil {
		return err
	}

	// 檢查 IP 限制（重用登入的 IP 限制）
	if err := s.checkIPLimit(ctx, ip); err != nil {
		return err
	}

	// 檢查該 email 的註冊嘗試次數（使用獨立的註冊計數器）
	key := fmt.Sprintf("register:email:%s", email)
	count, err := s.getCounter(ctx, key)
	if err != nil {
		return err
	}

	// 每個 email 每分鐘最多嘗試 3 次註冊
	if count >= 3 {
		return ErrRateLimitExceeded
	}

	return nil
}

// RecordRegisterAttempt 記錄註冊嘗試
func (s *RateLimitService) RecordRegisterAttempt(ctx context.Context, ip string, email string) error {
	// 記錄全局計數
	if err := s.incrementCounter(ctx, "login:global", s.config.GlobalWindow); err != nil {
		return err
	}

	// 記錄 IP 計數
	if err := s.incrementCounter(ctx, fmt.Sprintf("login:ip:%s", ip), s.config.IPWindow); err != nil {
		return err
	}

	// 記錄該 email 的註冊嘗試次數
	key := fmt.Sprintf("register:email:%s", email)
	if err := s.incrementCounter(ctx, key, s.config.UserWindow); err != nil {
		return err
	}

	return nil
}
