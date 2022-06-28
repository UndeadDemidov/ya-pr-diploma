package auth

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/app"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination=./mocks/mock_service.go . CredentialManager

type CredentialManager interface {
	AddNewUser(ctx context.Context, usr user.User, login, pword string) (err error)
	GetUser(ctx context.Context, login string) (usr user.User, err error)
	AuthenticateUser(ctx context.Context, login, pword string) (usr user.User, err error)
	// DisableUser(user user.User) error
}

var _ app.Authenticator = (*Service)(nil)

type Service struct {
	userSvc user.Registerer
	credMan CredentialManager
}

func NewService(userSvc user.Registerer, credMan CredentialManager) *Service {
	return &Service{
		userSvc: userSvc,
		credMan: credMan,
	}
}

func NewServiceWithDefaultCredMan(repo Repository, userSvc user.Registerer) *Service {
	return NewService(userSvc, NewManager(repo))
}

func (s Service) SignIn(ctx context.Context, login, pword string) error {
	// найти пользователя по логину - если есть, то занят
	_, err := s.credMan.GetUser(ctx, login)
	if err == nil {
		return errors2.ErrLoginIsInUseAlready
	}
	// если не занят, то создаем пустого пользователя и регистрируем его
	usr := user.NewUser()
	err = s.userSvc.RegisterNewUser(ctx, usr)
	if err != nil {
		return err
	}
	// создаем креды на пользователя
	err = s.credMan.AddNewUser(ctx, usr, login, pword)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) Login(ctx context.Context, login, pword string) (user user.User, err error) {
	return s.credMan.AuthenticateUser(ctx, login, pword)
}
