package balance

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
)

// var _ app.BalanceGetter = (*Service)(nil)

type Repository interface {
	Read(ctx context.Context, usr user.User) (Balance, error)
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
