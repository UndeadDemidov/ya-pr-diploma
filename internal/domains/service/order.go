package service

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/app"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/entity"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
)

var _ app.OrderProcessor = (*Order)(nil)

type Order struct {
}

func (o Order) Add(ctx context.Context, usr user.User, num string) error {
	// TODO implement me
	panic("implement me")
}

func (o Order) List(ctx context.Context, usr user.User) (ords []entity.Order, err error) {
	// TODO implement me
	panic("implement me")
}
