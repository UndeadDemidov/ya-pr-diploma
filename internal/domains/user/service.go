package user

import "context"

type Manager interface {
	RegisterNewUser(ctx context.Context, user User) error
}

type Repository interface {
	Create(ctx context.Context, user User) error
}

var _ Manager = (*Service)(nil)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterNewUser(ctx context.Context, user User) error {
	return s.repo.Create(ctx, user)
}
