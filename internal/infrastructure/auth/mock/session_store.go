package mock

import (
	"context"
	"sync"
	"time"

	"github.com/HiroLiang/goat-server/internal/domain/auth"
	"github.com/HiroLiang/goat-server/internal/infrastructure/auth/session"
)

type sessionEntry struct {
	session   *auth.Session
	expiresAt time.Time
}

type SessionStore struct {
	mu           sync.RWMutex
	sessions     map[string]*sessionEntry       // token -> session
	userSessions map[string]map[string]struct{} // userID -> set of tokens
}

func MockSessionStore() *SessionStore {
	return &SessionStore{
		sessions:     make(map[string]*sessionEntry),
		userSessions: make(map[string]map[string]struct{}),
	}
}

func (m *SessionStore) Get(_ context.Context, token string) (*auth.Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, exists := m.sessions[token]
	if !exists {
		return nil, ErrSessionNotFound
	}

	if time.Now().After(entry.expiresAt) {
		return nil, ErrSessionNotFound
	}

	return entry.session, nil
}

func (m *SessionStore) Set(_ context.Context, token string, sess *auth.Session, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[token] = &sessionEntry{
		session:   sess,
		expiresAt: time.Now().Add(ttl),
	}
	return nil
}

func (m *SessionStore) Delete(_ context.Context, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, token)
	return nil
}

func (m *SessionStore) Refresh(_ context.Context, token string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry, exists := m.sessions[token]
	if !exists {
		return ErrSessionNotFound
	}

	entry.expiresAt = time.Now().Add(ttl)
	return nil
}

func (m *SessionStore) AddUserSession(_ context.Context, userID, token string, _ time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.userSessions[userID] == nil {
		m.userSessions[userID] = make(map[string]struct{})
	}
	m.userSessions[userID][token] = struct{}{}
	return nil
}

func (m *SessionStore) RemoveUserSession(_ context.Context, userID, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if tokens, exists := m.userSessions[userID]; exists {
		delete(tokens, token)
	}
	return nil
}

func (m *SessionStore) ListUserSessions(_ context.Context, userID string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tokens, exists := m.userSessions[userID]
	if !exists {
		return []string{}, nil
	}

	result := make([]string, 0, len(tokens))
	for token := range tokens {
		result = append(result, token)
	}
	return result, nil
}

func (m *SessionStore) DeleteAllUserSessions(_ context.Context, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tokens, exists := m.userSessions[userID]
	if !exists {
		return nil
	}

	// 刪除所有相關的 session
	for token := range tokens {
		delete(m.sessions, token)
	}

	// 清空 user 的 token 列表
	delete(m.userSessions, userID)
	return nil
}

func (m *SessionStore) CleanExpiredUserSessions(_ context.Context, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tokens, exists := m.userSessions[userID]
	if !exists {
		return nil
	}

	now := time.Now()
	for token := range tokens {
		entry, exists := m.sessions[token]
		if !exists || now.After(entry.expiresAt) {
			delete(tokens, token)
			delete(m.sessions, token)
		}
	}
	return nil
}

var _ session.Store = (*SessionStore)(nil)
