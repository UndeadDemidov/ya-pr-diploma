package app

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/entity"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/primit"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination=./mocks/mock_gophermart.go . Authenticator,OrderProcessor,BalanceGetter,WithdrawalProcessor

type Authenticator interface {
	SignIn(ctx context.Context, login, pword string) error
	Login(ctx context.Context, login, pword string) (usr user.User, err error)
}

type OrderProcessor interface {
	Add(ctx context.Context, usr user.User, num string) error
	List(ctx context.Context, usr user.User) (ords []entity.Order, err error)
}

type BalanceGetter interface {
	Get(ctx context.Context, usr user.User) (bal entity.Balance, err error)
}

type WithdrawalProcessor interface {
	Add(ctx context.Context, usr user.User, num string, sum primit.Currency) error
	List(ctx context.Context, user2 user.User) (wtdrwls []entity.Withdrawal, err error)
}

type GopherMart struct {
	Authenticator
}

func NewGopherMart(auth Authenticator) *GopherMart {
	if auth == nil {
		panic("missing Authenticator, parameter must not be nil")
	}
	return &GopherMart{Authenticator: auth}
}
