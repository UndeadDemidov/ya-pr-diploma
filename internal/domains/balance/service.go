package balance

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/primit"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
)

// var (
// 	_ app.BalanceGetter = (*Service)(nil)
// 	_ app.WithdrawalProcessor = (*Service)(nil)
// )

type Repository interface {
	Read(context.Context, user.User) (Balance, error)
	CreateWithdrawal(context.Context, Withdrawal) error
	ListWithdrawals(context.Context, user.User) ([]Withdrawal, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("missing Repository, parameter must not be nil")
	}
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, usr user.User) (Balance, error) {
	return s.repo.Read(ctx, usr)
}

func (s *Service) Add(ctx context.Context, usr user.User, num primit.LuhnNumber, sum primit.Currency) error {
	wtdrwl, err := NewWithdrawal(usr, num, sum)
	if err != nil {
		return err
	}
	return s.repo.CreateWithdrawal(ctx, wtdrwl)
}

func (s *Service) List(ctx context.Context, usr user.User) ([]Withdrawal, error) {
	return s.repo.ListWithdrawals(ctx, usr)
}
