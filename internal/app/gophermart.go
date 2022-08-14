package app

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/balance"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/order"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/primit"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination=./mocks/mock_gophermart.go . Authenticator,OrderProcessor,BalanceGetter,WithdrawalProcessor

type Authenticator interface {
	SignIn(ctx context.Context, login, pword string) (user.User, error)
	Login(ctx context.Context, login, pword string) (user.User, error)
}

type OrderProcessor interface {
	// ToDo передавать сразу primit.LuhnNumber
	Add(ctx context.Context, usr user.User, num string) error
	List(context.Context, user.User) ([]order.Order, error)
	Close()
}

type BalanceGetter interface {
	Get(context.Context, user.User) (balance.Balance, error)
}

type WithdrawalProcessor interface {
	Add(context.Context, user.User, primit.LuhnNumber, primit.Currency) error
	List(context.Context, user.User) ([]balance.Withdrawal, error)
}

type GopherMart struct {
	Authenticator
	OrderProcessor
	BalanceGetter
}

func NewGopherMart(auth Authenticator, order OrderProcessor, bal BalanceGetter) *GopherMart {
	if auth == nil {
		panic("missing Authenticator, parameter must not be nil")
	}
	if order == nil {
		panic("missing OrderProcessor, parameter must not be nil")
	}
	if bal == nil {
		panic("missing BalanceGetter, parameter must not be nil")
	}
	return &GopherMart{
		Authenticator:  auth,
		OrderProcessor: order,
		BalanceGetter:  bal,
	}
}

func (m GopherMart) Close() {
	m.OrderProcessor.Close()
}
