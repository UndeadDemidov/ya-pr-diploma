package app

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
)

//go:generate mockgen -destination=./mocks/mock_gophermart.go . Authenticator

type Authenticator interface {
	SignIn(ctx context.Context, login, pword string) error
	Login(ctx context.Context, login, pword string) (user user.User, err error)
}

type GopherMart struct {
	Authenticator
}

func NewGopherMart(auth Authenticator) *GopherMart {
	return &GopherMart{Authenticator: auth}
}
