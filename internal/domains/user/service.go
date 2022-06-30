package user

import (
	"context"

	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination=./mocks/mock_service.go . Registerer,Repository

type Registerer interface {
	RegisterNewUser(ctx context.Context, user User) error
}

type Repository interface {
	Create(ctx context.Context, user User) error
}

var _ Registerer = (*Service)(nil)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("missing Repository, parameter must not be nil")
	}
	return &Service{repo: repo}
}

func (s *Service) RegisterNewUser(ctx context.Context, user User) error {
	return s.repo.Create(ctx, user)
}
