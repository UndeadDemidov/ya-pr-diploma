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

// Authenticator контракт для реализации регистрации и аутентификации пользователя.
type Authenticator interface {
	// SignIn регистрирует нового пользователя.
	SignIn(ctx context.Context, login, pword string) (user.User, error)
	// Login аутентифицирует существующего пользователя.
	Login(ctx context.Context, login, pword string) (user.User, error)
}

// OrderProcessor контракт для работы с заказами по начислениями.
type OrderProcessor interface {
	// Add добавляет для пользователя новый заказ для расчета начисленных баллов.
	Add(ctx context.Context, usr user.User, num string) error
	// List возвращает список заказов для пользователя ранее переданных для расчета баллов.
	List(context.Context, user.User) ([]order.Order, error)
	Close()
}

// BalanceGetter контракт для работы с балансом пользователя.
type BalanceGetter interface {
	// Get возвращает текущее состояние баланса пользователя.
	Get(context.Context, user.User) (balance.Balance, error)
}

// WithdrawalProcessor контракт для работы со списанием баллов в счет заказа.
type WithdrawalProcessor interface {
	// Add добавляет новый заказ для списания баллов.
	Add(context.Context, user.User, primit.LuhnNumber, primit.Currency) error
	// List возвращает ранее зарегистрированные списания пользователя.
	List(context.Context, user.User) ([]balance.Withdrawal, error)
}

// GopherMart является по сути приложением построенным через композицию контрактов.
type GopherMart struct {
	Authenticator
	OrderProcessor
	BalanceGetter
	WithdrawalProcessor
}

// NewGopherMart создает приложение GopherMart из составных частей композиции.
func NewGopherMart(
	auth Authenticator,
	order OrderProcessor,
	bal BalanceGetter,
	wtdrwl WithdrawalProcessor,
) *GopherMart {
	if auth == nil {
		panic("missing Authenticator, parameter must not be nil")
	}
	if order == nil {
		panic("missing OrderProcessor, parameter must not be nil")
	}
	if bal == nil {
		panic("missing BalanceGetter, parameter must not be nil")
	}
	if wtdrwl == nil {
		panic("missing WithdrawalProcessor, parameter must not be nil")
	}
	return &GopherMart{
		Authenticator:       auth,
		OrderProcessor:      order,
		BalanceGetter:       bal,
		WithdrawalProcessor: wtdrwl,
	}
}

// Close завершает приложение путем вызова вложенных Close методов.
func (m GopherMart) Close() {
	m.OrderProcessor.Close()
}
