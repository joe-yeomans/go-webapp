package store

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type User struct {
	Email           string
	IsEmailVerified bool
	GoogleID        string // Google OAuth ID
	GitHubID        string // GitHub OAuth ID
	Name            string // Full name from OAuth provider
	Picture         string // Profile picture URL from OAuth provider
	AuthProvider    string // "email", "google", or "github"
}

type LoginCode struct {
	Email     string
	Code      string
	ExpiresAt time.Time
}

type Session struct {
	UserEmail string
	Token     string
	ExpiresAt time.Time
}

type storeItem[T any] struct {
	Value T
	Lock  sync.RWMutex
}

type Store struct {
	users       storeItem[map[string]User]
	loginCodes  storeItem[map[string]LoginCode]
	sessions    storeItem[map[string]Session]
	sessionTime time.Duration
}

func NewStore(sessionTime time.Duration) *Store {
	return &Store{
		users:       storeItem[map[string]User]{Value: make(map[string]User)},
		loginCodes:  storeItem[map[string]LoginCode]{Value: make(map[string]LoginCode)},
		sessions:    storeItem[map[string]Session]{Value: make(map[string]Session)},
		sessionTime: sessionTime,
	}
}

func (s *Store) GetUser(email string) (User, bool) {
	s.users.Lock.RLock()
	defer s.users.Lock.RUnlock()
	user, ok := s.users.Value[email]
	return user, ok
}

func (s *Store) SetUser(email string, user User) {
	s.users.Lock.Lock()
	defer s.users.Lock.Unlock()
	s.users.Value[email] = user
}

func (s *Store) GetLoginCode(email string) (LoginCode, bool) {
	s.loginCodes.Lock.RLock()
	defer s.loginCodes.Lock.RUnlock()
	loginCode, ok := s.loginCodes.Value[email]
	return loginCode, ok
}

func (s *Store) SetLoginCode(email string, loginCode LoginCode) {
	s.loginCodes.Lock.Lock()
	defer s.loginCodes.Lock.Unlock()
	s.loginCodes.Value[email] = loginCode
}

func (s *Store) DeleteLoginCode(email string) {
	s.loginCodes.Lock.Lock()
	defer s.loginCodes.Lock.Unlock()
	delete(s.loginCodes.Value, email)
}

// Session management methods
func (s *Store) CreateSession(userEmail string) (string, error) {
	// Generate a secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	session := Session{
		UserEmail: userEmail,
		Token:     token,
		ExpiresAt: time.Now().Add(s.sessionTime), // 24 hours
	}

	s.sessions.Lock.Lock()
	defer s.sessions.Lock.Unlock()
	s.sessions.Value[token] = session

	return token, nil
}

func (s *Store) GetSession(token string) (Session, bool) {
	s.sessions.Lock.RLock()
	defer s.sessions.Lock.RUnlock()
	session, ok := s.sessions.Value[token]
	return session, ok
}

func (s *Store) DeleteSession(token string) {
	s.sessions.Lock.Lock()
	defer s.sessions.Lock.Unlock()
	delete(s.sessions.Value, token)
}

func (s *Store) IsSessionValid(token string) bool {
	session, exists := s.GetSession(token)
	if !exists {
		return false
	}
	return time.Now().Before(session.ExpiresAt)
}

func (s *Store) CleanupExpiredSessions() {
	s.sessions.Lock.Lock()
	defer s.sessions.Lock.Unlock()

	now := time.Now()
	for token, session := range s.sessions.Value {
		if now.After(session.ExpiresAt) {
			delete(s.sessions.Value, token)
		}
	}
}

func (s *Store) GetSessionTime() time.Duration {
	return s.sessionTime
}
