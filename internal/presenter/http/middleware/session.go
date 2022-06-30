package middleware

import (
	"time"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/app"
	"github.com/google/uuid"
)

const sessionTTL = 120 * time.Second

type Session struct {
	userRef app.Referencer
	expiry  time.Time
}

func NewSession(userRef app.Referencer, ttl time.Duration) Session {
	if userRef == nil {
		panic("missing app.Referencer, parameter must not be nil")
	}
	s := Session{userRef: userRef}
	s.Refresh(ttl)
	return s
}

func (s *Session) IsExpired() bool {
	return s.expiry.Before(time.Now())
}

func (s *Session) Refresh(ttl time.Duration) {
	s.expiry = time.Now().Add(ttl)
}

func (s *Session) GetReference() string {
	return s.userRef.Reference()
}

type SessionToken string

func NewSessionToken() SessionToken {
	return SessionToken(uuid.New().String())
}

type Sessions struct {
	ttl   time.Duration
	store map[SessionToken]Session
}

func NewSessions(ttl time.Duration) *Sessions {
	return &Sessions{
		ttl:   ttl,
		store: make(map[SessionToken]Session, 8),
	}
}

func NewDefaultSessions() *Sessions {
	return NewSessions(sessionTTL)
}

func (s *Sessions) AddNewSession(userRef app.Referencer) SessionToken {
	if userRef == nil {
		panic("missing app.Referencer, parameter must not be nil")
	}
	token := NewSessionToken()
	s.store[token] = NewSession(userRef, s.ttl)
	return token
}

func (s *Sessions) IsExpired(token SessionToken) bool {
	if session, ok := s.store[token]; ok {
		if expired := session.IsExpired(); expired {
			delete(s.store, token)
			return true
		}
		session.Refresh(s.ttl)
		return false
	}
	return true
}

func (s *Sessions) GetReference(token SessionToken) string {
	if session, ok := s.store[token]; ok {
		return session.userRef.Reference()
	}
	return ""
}
