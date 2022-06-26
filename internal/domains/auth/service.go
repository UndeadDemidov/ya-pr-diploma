package auth

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/app"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
)

//go:generate mockgen -destination=./mocks/mock_service.go . CredentialManager

type CredentialManager interface {
	AddNewUser(ctx context.Context, user user.User, login, pword string) (err error)
	GetUser(ctx context.Context, login, pword string) (user user.User, err error)
	// DisableUser(user user.User) error
}

var _ app.Authenticator = (*Service)(nil)

type Service struct {
	// userSvc
	credMan CredentialManager
}

func NewService(credMan CredentialManager) *Service {
	return &Service{credMan: credMan}
}

func NewServiceWithDefaultCredMan() *Service {
	return NewService(&Credentials{})
}

func (s Service) SignIn(ctx context.Context, login, pword string) error {
	// найти пользователя по логину - если есть, то занят
	_, err := s.credMan.GetUser(ctx, login, pword)
	if err == nil {
		return errors2.ErrLoginIsInUseAlready
	}
	// если не занят, то создаем пустого пользователя
	err = s.credMan.AddNewUser(ctx, user.NewUser(), login, pword)
	if err != nil {
		return err
	}
	// создаем креды на пользователя
	return errors2.ErrLoginIsInUseAlready
}

func (s Service) Login(ctx context.Context, login, pword string) (user user.User, err error) {
	// найти пользователя,
	// если не найден - то ошибка
	// нужно ли найти пользователя или по сессии его надо находить?
	// TODO implement me
	panic("implement me")
}
