package order

import (
	"context"
	"strconv"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/primit"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
)

// var _ app.OrderProcessor = (*Service)(nil)

type Repository interface {
	Create(ctx context.Context, ord Order) error
	ListByUser(ctx context.Context, usr user.User) ([]Order, error)
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

func (s Service) Add(ctx context.Context, usr user.User, num string) error {
	num64, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return err
	}
	lnum := primit.LuhnNumber(num64)
	ord, err := NewOrder(usr, lnum)
	if err != nil {
		return err
	}
	err = s.repo.Create(ctx, ord)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) List(ctx context.Context, usr user.User) (ords []Order, err error) {
	return s.repo.ListByUser(ctx, usr)
}
