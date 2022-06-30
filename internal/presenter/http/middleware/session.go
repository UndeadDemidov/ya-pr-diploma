package middleware

import (
	"time"

	"github.com/google/uuid"
)

const sessionTTL = 120 * time.Second

type Referencer interface {
	Reference() string
}

type session struct {
	userRef Referencer
	expiry  time.Time
}

func newSession(userRef Referencer, ttl time.Duration) session {
	if userRef == nil {
		panic("missing app.Referencer, parameter must not be nil")
	}
	s := session{userRef: userRef}
	s.refresh(ttl)
	return s
}

func (s *session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func (s *session) refresh(ttl time.Duration) {
	s.expiry = time.Now().Add(ttl)
}

func (s *session) getReference() string {
	return s.userRef.Reference()
}

type SessionToken string

func NewSessionToken() SessionToken {
	return SessionToken(uuid.New().String())
}

type Sessions struct {
	ttl   time.Duration
	store map[SessionToken]session
}

func NewSessions(ttl time.Duration) *Sessions {
	return &Sessions{
		ttl:   ttl,
		store: make(map[SessionToken]session, 8),
	}
}

func NewDefaultSessions() *Sessions {
	return NewSessions(sessionTTL)
}

func (s *Sessions) AddNewSession(userRef Referencer) SessionToken {
	if userRef == nil {
		panic("missing app.Referencer, parameter must not be nil")
	}
	token := NewSessionToken()
	s.store[token] = newSession(userRef, s.ttl)
	return token
}

func (s *Sessions) IsExpired(token SessionToken) bool {
	if session, ok := s.store[token]; ok {
		if expired := session.isExpired(); expired {
			delete(s.store, token)
			return true
		}
		session.refresh(s.ttl)
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
